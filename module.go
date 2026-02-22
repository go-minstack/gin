package gin

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("gin",
		fx.Provide(NewServer),
	)
}
