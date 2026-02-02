package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"polymarket-cli/internal/client"
	"polymarket-cli/internal/config"
)

var (
	market        []string
	eventID       []int
	sizeThreshold float64
	redeemable    bool
	mergeable     bool
	limit         int
	offset        int
	sortBy        string
	sortDirection string
	title         string
)

var positionsCmd = &cobra.Command{
	Use:   "positions [user-address]",
	Short: "Get current positions for a user",
	Long:  `Returns positions filtered by user and optional filters.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: user address is required")
			return
		}

		userAddr := args[0]

		positions, err := fetchPositions(userAddr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		jsonData, err := json.MarshalIndent(positions, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting output: %v\n", err)
			return
		}

		fmt.Println(string(jsonData))
	},
}

func init() {
	rootCmd.AddCommand(positionsCmd)

	positionsCmd.Flags().StringSliceVar(&market, "market", []string{}, "Comma-separated list of condition IDs")
	positionsCmd.Flags().IntSliceVar(&eventID, "event-id", []int{}, "Comma-separated list of event IDs")
	positionsCmd.Flags().Float64Var(&sizeThreshold, "size-threshold", 1, "Minimum size threshold")
	positionsCmd.Flags().BoolVar(&redeemable, "redeemable", false, "Filter redeemable positions")
	positionsCmd.Flags().BoolVar(&mergeable, "mergeable", false, "Filter mergeable positions")
	positionsCmd.Flags().IntVar(&limit, "limit", 100, "Limit results (0-500)")
	positionsCmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination (0-10000)")
	positionsCmd.Flags().StringVar(&sortBy, "sort-by", "TOKENS", "Sort by (CURRENT, INITIAL, TOKENS, CASHPNL, PERCENTPNL, TITLE, RESOLVING, PRICE, AVGPRICE)")
	positionsCmd.Flags().StringVar(&sortDirection, "sort-direction", "DESC", "Sort direction (ASC, DESC)")
	positionsCmd.Flags().StringVar(&title, "title", "", "Filter by title")
}

type Position struct {
	ProxyWallet        string  `json:"proxyWallet"`
	Asset              string  `json:"asset"`
	ConditionID        string  `json:"conditionId"`
	Size               float64 `json:"size"`
	AvgPrice           float64 `json:"avgPrice"`
	InitialValue       float64 `json:"initialValue"`
	CurrentValue       float64 `json:"currentValue"`
	CashPnl            float64 `json:"cashPnl"`
	PercentPnl         float64 `json:"percentPnl"`
	TotalBought        float64 `json:"totalBought"`
	RealizedPnl        float64 `json:"realizedPnl"`
	PercentRealizedPnl float64 `json:"percentRealizedPnl"`
	CurPrice           float64 `json:"curPrice"`
	Redeemable         bool    `json:"redeemable"`
	Mergeable          bool    `json:"mergeable"`
	Title              string  `json:"title"`
	Slug               string  `json:"slug"`
	Icon               string  `json:"icon"`
	EventSlug          string  `json:"eventSlug"`
	Outcome            string  `json:"outcome"`
	OutcomeIndex       int     `json:"outcomeIndex"`
	OppositeOutcome    string  `json:"oppositeOutcome"`
	OppositeAsset      string  `json:"oppositeAsset"`
	EndDate            string  `json:"endDate"`
	NegativeRisk       bool    `json:"negativeRisk"`
}

func fetchPositions(userAddr string) ([]Position, error) {
	httpClient := client.NewHTTPClient(config.AppCfg.DataAPIBaseURL)

	query := url.Values{}
	query.Set("user", userAddr)

	if len(market) > 0 {
		for _, m := range market {
			query.Add("market", m)
		}
	}

	if len(eventID) > 0 {
		for _, id := range eventID {
			query.Add("eventId", fmt.Sprintf("%d", id))
		}
	}

	query.Set("sizeThreshold", fmt.Sprintf("%.0f", sizeThreshold))

	if redeemable {
		query.Set("redeemable", "true")
	} else {
		query.Set("redeemable", "false")
	}

	if mergeable {
		query.Set("mergeable", "true")
	} else {
		query.Set("mergeable", "false")
	}

	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))

	if sortBy != "" {
		query.Set("sortBy", sortBy)
	}

	if sortDirection != "" {
		query.Set("sortDirection", sortDirection)
	}

	if title != "" {
		query.Set("title", title)
	}

	var positions []Position
	if err := httpClient.GetJSONWithMultipleValues("/positions", query, &positions); err != nil {
		return nil, err
	}

	return positions, nil
}
