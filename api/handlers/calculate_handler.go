package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/romutchio/crypto-calculator/internal/entity"
	"github.com/romutchio/crypto-calculator/internal/usecase"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetPair(ctx context.Context, from string, to string) (*entity.Pair, error)
}

type CalculateResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
}

// CalculateHandler
//
//	@Summary		CalculateHandler crypto/fiat exchange
//	@Param			from	query		string		    true	"from"
//	@Param			to		query		string		    true	"to"
//	@Param			amount	query		number		    true	"amount"
//	@Success		200		{string}	string			"ok"
//	@Router			/api/v1/calculate [get]
func CalculateHandler(ctx context.Context, uc usecase.CalculateUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")
		amount := c.QueryFloat("amount")
		resp, err := uc.Calculate(ctx, &usecase.CalculateRequest{
			From:   from,
			To:     to,
			Amount: amount,
		})
		if err != nil {
			log.Err(err).Msg("uc.CalculateHandler failed")
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.JSON(&CalculateResponse{
			From:   from,
			To:     to,
			Rate:   resp.Rate,
			Amount: amount,
			Result: resp.Amount,
		})
	}
}
