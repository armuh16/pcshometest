package module

import (
	//Route
	authRoute "pcstakehometest/module/auth/route"
	productRoute "pcstakehometest/module/product/route"
	transactionRoute "pcstakehometest/module/transaction/route"

	//Logic
	authLogic "pcstakehometest/module/auth/logic"
	productLogic "pcstakehometest/module/product/logic"
	transactionLogic "pcstakehometest/module/transaction/logic"
	userLogic "pcstakehometest/module/user/logic"

	//Repository
	productRepository "pcstakehometest/module/product/repository"
	transactionRepository "pcstakehometest/module/transaction/repository"
	userRepository "pcstakehometest/module/user/repository"

	"go.uber.org/fx"
)

// Register Route
var BundleRoute = fx.Options(
	fx.Invoke(transactionRoute.NewRoute),
	fx.Invoke(productRoute.NewRoute),
	fx.Invoke(authRoute.NewRoute),
)

// Register logic
var BundleLogic = fx.Options(
	fx.Provide(userLogic.NewLogic),
	fx.Provide(transactionLogic.NewLogic),
	fx.Provide(productLogic.NewLogic),
	fx.Provide(authLogic.NewLogic),
)

// Register Repository
var BundleRepository = fx.Options(
	fx.Provide(userRepository.NewRepository),
	fx.Provide(transactionRepository.NewRepository),
	fx.Provide(productRepository.NewRepository),
)
