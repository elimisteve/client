// Auto-generated by avdl-compiler v1.3.22 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/keybase1/avatars.avdl

package keybase1

import (
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type AvatarUrl string

func (o AvatarUrl) DeepCopy() AvatarUrl {
	return o
}

type AvatarFormat string

func (o AvatarFormat) DeepCopy() AvatarFormat {
	return o
}

type LoadUserAvatarsRes struct {
	Picmap map[string]map[AvatarFormat]AvatarUrl `codec:"picmap" json:"picmap"`
}

func (o LoadUserAvatarsRes) DeepCopy() LoadUserAvatarsRes {
	return LoadUserAvatarsRes{
		Picmap: (func(x map[string]map[AvatarFormat]AvatarUrl) map[string]map[AvatarFormat]AvatarUrl {
			if x == nil {
				return nil
			}
			ret := make(map[string]map[AvatarFormat]AvatarUrl)
			for k, v := range x {
				kCopy := k
				vCopy := (func(x map[AvatarFormat]AvatarUrl) map[AvatarFormat]AvatarUrl {
					if x == nil {
						return nil
					}
					ret := make(map[AvatarFormat]AvatarUrl)
					for k, v := range x {
						kCopy := k.DeepCopy()
						vCopy := v.DeepCopy()
						ret[kCopy] = vCopy
					}
					return ret
				})(v)
				ret[kCopy] = vCopy
			}
			return ret
		})(o.Picmap),
	}
}

type LoadUserAvatarsArg struct {
	Usernames []string       `codec:"usernames" json:"usernames"`
	Formats   []AvatarFormat `codec:"formats" json:"formats"`
}

type AvatarsInterface interface {
	LoadUserAvatars(context.Context, LoadUserAvatarsArg) (LoadUserAvatarsRes, error)
}

func AvatarsProtocol(i AvatarsInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "keybase.1.avatars",
		Methods: map[string]rpc.ServeHandlerDescription{
			"loadUserAvatars": {
				MakeArg: func() interface{} {
					ret := make([]LoadUserAvatarsArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadUserAvatarsArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadUserAvatarsArg)(nil), args)
						return
					}
					ret, err = i.LoadUserAvatars(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type AvatarsClient struct {
	Cli rpc.GenericClient
}

func (c AvatarsClient) LoadUserAvatars(ctx context.Context, __arg LoadUserAvatarsArg) (res LoadUserAvatarsRes, err error) {
	err = c.Cli.Call(ctx, "keybase.1.avatars.loadUserAvatars", []interface{}{__arg}, &res)
	return
}
