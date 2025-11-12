package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RazorpayClient struct {
	keyID     string
	keySecret string
	baseURL   string
}

func NewRazorpayClient(keyID, keySecret string) *RazorpayClient {
	return &RazorpayClient{
		keyID:     keyID,
		keySecret: keySecret,
		baseURL:   "https://api.razorpay.com/v1",
	}
}

// CreateOrder creates a new Razorpay order
func (r *RazorpayClient) CreateOrder(amount int, currency, receipt string) (string, error) {
	data := map[string]interface{}{
		"amount":   amount,
		"currency": currency,
		"receipt":  receipt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", r.baseURL+"/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(r.keyID, r.keySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create Razorpay order: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("razorpay API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	orderID, ok := result["id"].(string)
	if !ok {
		return "", errors.New("failed to get order ID from Razorpay response")
	}

	return orderID, nil
}

// VerifyPaymentSignature verifies the Razorpay payment signature
func (r *RazorpayClient) VerifyPaymentSignature(orderID, paymentID, signature string) error {
	// Create the message to verify
	message := orderID + "|" + paymentID

	// Create HMAC SHA256 hash
	h := hmac.New(sha256.New, []byte(r.keySecret))
	h.Write([]byte(message))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare signatures
	if expectedSignature != signature {
		return errors.New("invalid payment signature")
	}

	return nil
}

// GetPaymentDetails retrieves payment details from Razorpay
func (r *RazorpayClient) GetPaymentDetails(paymentID string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", r.baseURL+"/payments/"+paymentID, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(r.keyID, r.keySecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment details: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// RefundPayment initiates a refund for a payment
func (r *RazorpayClient) RefundPayment(paymentID string, amount int) (string, error) {
	data := map[string]interface{}{
		"amount": amount,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", r.baseURL+"/payments/"+paymentID+"/refund", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(r.keyID, r.keySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to refund payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	refundID, ok := result["id"].(string)
	if !ok {
		return "", errors.New("failed to get refund ID from Razorpay response")
	}

	return refundID, nil
}
