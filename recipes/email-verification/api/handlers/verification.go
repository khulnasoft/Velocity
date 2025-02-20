package handlers

import (
	"email-verification/application"

	"go.khulnasoft.com/velocity"
)

type VerificationHandler struct {
	verificationService *application.VerificationService
}

func NewVerificationHandler(service *application.VerificationService) *VerificationHandler {
	return &VerificationHandler{verificationService: service}
}

func (h *VerificationHandler) SendVerification(c *velocity.Ctx) error {
	email := c.Params("email")
	if err := h.verificationService.SendVerification(email); err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(velocity.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(velocity.Map{"message": "Verification code sent"})
}

func (h *VerificationHandler) CheckVerification(c *velocity.Ctx) error {
	email := c.Params("email")
	code := c.Params("code")

	if err := h.verificationService.VerifyCode(email, code); err != nil {
		return c.Status(velocity.StatusBadRequest).JSON(velocity.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(velocity.Map{"message": "Code verified successfully"})
}
