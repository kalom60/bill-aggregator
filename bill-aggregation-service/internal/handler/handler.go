package handler

import (
	"bill-aggregation-service/internal/database"
	"bill-aggregation-service/internal/grpc/clients"
	"bill-aggregation-service/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetBills(c *gin.Context)
	GetBillsByProvider(c *gin.Context)
	RefreshBills(c *gin.Context)
}

type handler struct {
	cache              database.Service
	gRPCAccountClient  *clients.AccountClient
	gRPCProviderClient *clients.ProviderClient
}

func NewHandler(cache database.Service, gRPCAccountClient *clients.AccountClient, gRPCProviderClient *clients.ProviderClient) Handler {
	return &handler{
		cache:              cache,
		gRPCAccountClient:  gRPCAccountClient,
		gRPCProviderClient: gRPCProviderClient,
	}
}

func (h *handler) fetchBills(c *gin.Context, userID string) ([]models.Bill, error) {
	accounts, err := h.gRPCAccountClient.GetLinkedAccounts(c.Request.Context(), userID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	billChan := make(chan models.Bill, len(accounts))
	errChan := make(chan error, len(accounts))

	for _, account := range accounts {
		wg.Add(1)
		go func(a clients.Account) {
			defer wg.Done()

			provider, err := h.gRPCProviderClient.GetProvider(c.Request.Context(), a.ProviderId)
			if err != nil {
				errChan <- err
				return
			}

			bills, err := h.fetchProviderBills(provider)
			if err != nil {
				errChan <- err
				return
			}

			billChan <- models.Bill{
				UserID:       userID,
				ProviderID:   provider.Id,
				ProviderName: provider.Name,
				Items:        bills,
			}
		}(account)
	}

	go func() {
		wg.Wait()
		close(billChan)
		close(errChan)
	}()

	var allBills []models.Bill
	for bill := range billChan {
		allBills = append(allBills, bill)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return allBills, nil
}

func (h *handler) fetchProviderBills(provider *clients.Provider) ([]models.BillItem, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", provider.ApiURL, provider.ApiKey), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch bills, status: %d", response.StatusCode)
	}

	var billResponse models.BillResponse
	if err := json.NewDecoder(response.Body).Decode(&billResponse); err != nil {
		return nil, err
	}

	return billResponse.Bills, nil
}

func (h *handler) GetBills(c *gin.Context) {
	userID := c.Param("user_id")

	if bills, err := h.cache.Get(c.Request.Context(), userID); err == nil {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": bills})
		return
	}

	bills, err := h.fetchBills(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed while fetching bills"})
		return
	}

	if err := h.cache.Set(c.Request.Context(), userID, bills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed while fetching bills"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": bills})
}

func (h *handler) GetBillsByProvider(c *gin.Context) {
	userID := c.Param("user_id")
	providerName := c.Query("name")

	if providerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Provider name is required"})
		return
	}

	if cachedBills, err := h.cache.Get(c.Request.Context(), userID); err == nil {
		filteredBills := filterBillsByProvider(cachedBills, providerName)
		if len(filteredBills) > 0 {
			c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": filteredBills})
			return
		}
	}

	bills, err := h.fetchBills(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed while fetching bills"})
		return
	}

	if err := h.cache.Set(c.Request.Context(), userID, bills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed while fetching bills"})
		return
	}

	filteredBills := filterBillsByProvider(bills, providerName)
	if len(filteredBills) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "No bills found for the specified provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": filteredBills})
}

func (h *handler) RefreshBills(c *gin.Context) {
	userID := c.Param("user_id")

	bills, err := h.fetchBills(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed while refreshing bills"})
		return
	}

	if err := h.cache.Set(c.Request.Context(), userID, bills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed while fetching bills"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Bills refreshed successfully"})
}

func filterBillsByProvider(bills []models.Bill, providerName string) []models.Bill {
	var filtered []models.Bill
	for _, bill := range bills {
		if bill.ProviderName == providerName {
			filtered = append(filtered, bill)
		}
	}
	return filtered
}
