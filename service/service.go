package service

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewLoginService),
	fx.Provide(NewJWTAuthService),
	fx.Provide(NewBucketService),
)
