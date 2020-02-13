// +build !confonly

package router

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/routing"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		r := new(Router)
		if err := core.RequireFeatures(ctx, func(ohm outbound.Manager) error {
			return r.Init(config.(*Config), ohm)
		}); err != nil {
			return nil, err
		}
		return r, nil
	}))
}

// Router is an implementation of routing.Router.
type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []*Rule
}

// Init initializes the Router.
func (r *Router) Init(config *Config, ohm outbound.Manager) error {
	r.domainStrategy = config.DomainStrategy

	r.rules = make([]*Rule, 0, len(config.Rule))
	for _, rule := range config.Rule {
		cond, err := rule.BuildCondition()
		if err != nil {
			return err
		}
		rr := &Rule{
			Condition: cond,
			Tag:       rule.GetTag(),
		}
		r.rules = append(r.rules, rr)
	}

	return nil
}

func (r *Router) PickRoute(ctx context.Context) (string, error) {
	rule, err := r.pickRouteInternal(ctx)
	if err != nil {
		return "", err
	}
	return rule.GetTag()
}

func isDomainOutbound(outbound *session.Outbound) bool {
	return outbound != nil && outbound.Target.IsValid() && outbound.Target.Address.Family().IsDomain()
}

// PickRoute implements routing.Router.
func (r *Router) pickRouteInternal(ctx context.Context) (*Rule, error) {
	sessionContext := &Context{
		Inbound:  session.InboundFromContext(ctx),
		Outbound: session.OutboundFromContext(ctx),
		Content:  session.ContentFromContext(ctx),
	}

	for _, rule := range r.rules {
		if rule.Apply(sessionContext) {
			return rule, nil
		}
	}

	return nil, common.ErrNoClue
}

// Start implements common.Runnable.
func (*Router) Start() error {
	return nil
}

// Close implements common.Closable.
func (*Router) Close() error {
	return nil
}

// Type implement common.HasType.
func (*Router) Type() interface{} {
	return routing.RouterType()
}

type Context struct {
	Inbound  *session.Inbound
	Outbound *session.Outbound
	Content  *session.Content
}

func (c *Context) GetTargetIPs() []net.IP {
	if c.Outbound == nil || !c.Outbound.Target.IsValid() {
		return nil
	}

	if c.Outbound.Target.Address.Family().IsIP() {
		return []net.IP{c.Outbound.Target.Address.IP()}
	}

	if len(c.Outbound.ResolvedIPs) > 0 {
		return c.Outbound.ResolvedIPs
	}

	return nil
}
