package logos

import "context"

type ctxLoggerKey struct{}

func ToCtx(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func FromCtx(ctx context.Context) Logger {
	return ctx.Value(ctxLoggerKey{}).(Logger)
}
