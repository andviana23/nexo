package middleware

import (
	"github.com/labstack/echo/v4"
)

// UnitMiddleware ensures that a Unit ID is present in the context.
// It relies on JWTMiddleware to have already extracted it (from Claim or Header).
func UnitMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			unitID := GetUnitID(c)
			if unitID == "" {
				// Tenta pegar do header caso o JWTMiddleware falhe ou não seja usado
				unitID = c.Request().Header.Get("X-Unit-ID")
				if unitID == "" {
					return echo.NewHTTPError(400, "X-Unit-ID header is required")
				}
				c.Set("unit_id", unitID)
			}
			return next(c)
		}
	}
}
