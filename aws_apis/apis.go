package aws_apis

import (
	"context"
	"github.com/gin-gonic/gin"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/elmntri/zeitgeber-aws-modules/bucket_connector"
	"github.com/elmntri/zeitgeber-common-modules/http_server"
)

type APIs struct {
	params Params
	logger *zap.Logger
	router *gin.RouterGroup
	scope  string
}

type Params struct {
	fx.In

	Lifecycle  		fx.Lifecycle
	Logger     		*zap.Logger
	HTTPServer 		*http_server.HTTPServer
	BucketConnector	*bucket_connector.BucketConnector
}

func Module(scope string) fx.Option {
	var a *APIs
	return fx.Options(
		fx.Provide(func(p Params) *APIs {
			a := &APIs{
				params: p,
				logger: p.Logger.Named(scope),
				scope:  scope,
			}
			return a
		}),
		fx.Populate(&a),
		fx.Invoke(func(p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: a.onStart,
					OnStop:  a.onStop,
				},
			)
		}),
	)
}

func (a *APIs) onStart(ctx context.Context) error {

	a.logger.Info("Starting APIs")

	// Router
	a.router = a.params.HTTPServer.GetRouter().Group("apis/v1/aws")
	a.router.GET("/list_buckets", a.listBuckets)
	a.router.POST("/upload_to_bucket", a.uploadFile)

	return nil
}

func (a *APIs) onStop(ctx context.Context) error {
	a.logger.Info("Stopped APIs")

	return nil
}



