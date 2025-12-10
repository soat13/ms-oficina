package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/soat13/fase-1-oficina/pkg/pagination"
	string "github.com/soat13/fase-1-oficina/pkg/utils/helpers/string"
)

func NewPagination(ctx *fiber.Ctx, defaultLimit, defaultOffset int) *pagination.Pagination {
	limit := string.StringToIntOrDefault(ctx.Query("limit"), defaultLimit)
	offset := string.StringToIntOrDefault(ctx.Query("offset"), defaultOffset)
	return pagination.New(limit, offset)
}
