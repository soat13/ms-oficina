package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/soat13/fase-1-oficina/pkg/pagination"
	stringHelper "github.com/soat13/fase-1-oficina/pkg/utils/helpers/string"
)

func NewPagination(ctx *fiber.Ctx, defaultLimit, defaultOffset int) *pagination.Pagination {
	limit := stringHelper.StringToIntOrDefault(ctx.Query("limit"), defaultLimit)
	offset := stringHelper.StringToIntOrDefault(ctx.Query("offset"), defaultOffset)
	return pagination.New(limit, offset)
}
