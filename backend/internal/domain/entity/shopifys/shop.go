// Package shopifys
package shopifys

import "time"

// Shop 店铺信息
type Shop struct {
	AccountOwner struct {
		Name string `json:"name"`
	} `json:"accountOwner"`
	Alerts []struct {
		Action struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"action"`
		Description string `json:"description"`
	} `json:"alerts"`
	BillingAddress               Address  `json:"billingAddress"`
	CheckoutApiSupported         bool     `json:"checkoutApiSupported"`
	ContactEmail                 string   `json:"contactEmail"`
	CreatedAt                    string   `json:"createdAt"`
	CurrencyCode                 string   `json:"currencyCode"`
	CurrencyFormats              Formats  `json:"currencyFormats"`
	CustomerAccounts             string   `json:"customerAccounts"`
	Description                  string   `json:"description"`
	Email                        string   `json:"email"`
	EnabledPresentmentCurrencies []string `json:"enabledPresentmentCurrencies"`
	FulfillmentServices          []struct {
		Handle      string `json:"handle"`
		ServiceName string `json:"serviceName"`
	} `json:"fulfillmentServices"`
	IanaTimezone                         string `json:"ianaTimezone"`
	ID                                   string `json:"id"`
	MarketingSmsConsentEnabledAtCheckout bool   `json:"marketingSmsConsentEnabledAtCheckout"`
	MyshopifyDomain                      string `json:"myshopifyDomain"`
	Name                                 string `json:"name"`
	PaymentSettings                      struct {
		SupportedDigitalWallets []string `json:"supportedDigitalWallets"`
	} `json:"paymentSettings"`
	Plan          Plan `json:"plan"`
	PrimaryDomain struct {
		Host string `json:"host"`
		ID   string `json:"id"`
	} `json:"primaryDomain"`
	ProductTypes struct {
		Edges []struct {
			Node string `json:"node"`
		} `json:"edges"`
	} `json:"productTypes"`
	SetupRequired            bool     `json:"setupRequired"`
	ShipsToCountries         []string `json:"shipsToCountries"`
	TaxesIncluded            bool     `json:"taxesIncluded"`
	TaxShipping              bool     `json:"taxShipping"`
	TimezoneAbbreviation     string   `json:"timezoneAbbreviation"`
	TimezoneOffsetMinutes    int      `json:"timezoneOffsetMinutes,omitempty"`
	TransactionalSmsDisabled bool     `json:"transactionalSmsDisabled"`
	UpdatedAt                string   `json:"updatedAt"`
	URL                      string   `json:"url"`
	WeightUnit               string   `json:"weightUnit"`
}

// ShopResponse 店铺信息响应
type ShopResponse struct {
	Shop Shop `json:"shop"`
}

// Domain 域名信息
type Domain struct {
	URL  string `json:"url"`
	Host string `json:"host"`
	ID   string `json:"id,omitempty"`
}

// Plan 店铺套餐
type Plan struct {
	DisplayName        string `json:"displayName"`
	PartnerDevelopment bool   `json:"partnerDevelopment"`
	ShopifyPlus        bool   `json:"shopifyPlus"`
}

// Formats 货币格式
type Formats struct {
	MoneyFormat                     string `json:"moneyFormat"`
	MoneyInEmailsFormat             string `json:"moneyInEmailsFormat"`
	MoneyWithCurrencyFormat         string `json:"moneyWithCurrencyFormat"`
	MoneyWithCurrencyInEmailsFormat string `json:"moneyWithCurrencyInEmailsFormat"`
}

// Address 地址信息
type Address struct {
	Address1      string  `json:"address1"`
	Address2      string  `json:"address2"`
	City          string  `json:"city"`
	Company       string  `json:"company"`
	Country       string  `json:"country"`
	CountryCodeV2 string  `json:"countryCodeV2"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Phone         string  `json:"phone"`
	Province      string  `json:"province"`
	ProvinceCode  string  `json:"provinceCode"`
	Zip           string  `json:"zip"`
}

// App 应用信息
type App struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ShopBillingAddressInput 账单地址输入
type ShopBillingAddressInput struct {
	Address1     string `json:"address1"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city"`
	Company      string `json:"company,omitempty"`
	CountryCode  string `json:"countryCode"`
	Phone        string `json:"phone,omitempty"`
	ProvinceCode string `json:"provinceCode,omitempty"`
	Zip          string `json:"zip"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
}

// ShopSettingsInput 店铺设置输入
type ShopSettingsInput struct {
	Name                  string   `json:"name,omitempty"`
	ContactEmail          string   `json:"contactEmail,omitempty"`
	CustomerEmail         string   `json:"customerEmail,omitempty"`
	CheckoutApiSupported  bool     `json:"checkoutApiSupported,omitempty"`
	EnabledFeatures       []string `json:"enabledFeatures,omitempty"`
	Password              string   `json:"password,omitempty"`
	TimezoneAbbreviation  string   `json:"timezoneAbbreviation,omitempty"`
	TimezoneOffset        string   `json:"timezoneOffset,omitempty"`
	TimezoneIdentifier    string   `json:"timezoneIdentifier,omitempty"`
	TimezoneOffsetMinutes int      `json:"timezoneOffsetMinutes,omitempty"`
	UnitSystem            string   `json:"unitSystem,omitempty"`
	WeightsUnit           string   `json:"weightsUnit,omitempty"`
}

// ShopPoliciesResponse 店铺政策响应
type ShopPoliciesResponse struct {
	Shop struct {
		ID                 string `json:"id"`
		PrivacyPolicy      Policy `json:"privacyPolicy"`
		RefundPolicy       Policy `json:"refundPolicy"`
		TermsOfService     Policy `json:"termsOfService"`
		ShippingPolicy     Policy `json:"shippingPolicy"`
		SubscriptionPolicy Policy `json:"subscriptionPolicy"`
	} `json:"shop"`
}

// Policy 政策信息
type Policy struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

// ShopLocalesResponse 店铺语言响应
type ShopLocalesResponse struct {
	Shop struct {
		ID                    string `json:"id"`
		PrimaryDomain         Domain `json:"primaryDomain"`
		PrimaryLocale         string `json:"primaryLocale"`
		URL                   string `json:"url"`
		TranslatableContentV2 struct {
			Edges []struct {
				Node struct {
					Key                 string `json:"key"`
					Value               string `json:"value"`
					Digest              string `json:"digest"`
					TranslatableContent struct {
						Key    string `json:"key"`
						Value  string `json:"value"`
						Digest string `json:"digest"`
					} `json:"translatableContent"`
					Translations []struct {
						Locale string `json:"locale"`
						Key    string `json:"key"`
						Value  string `json:"value"`
					} `json:"translations"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"translatableContentV2"`
	} `json:"shop"`
}

// ShopInfo 响应的结构体
type ShopInfoResponse struct {
	Shop struct {
		Name            string `json:"name"`
		MyShopifyDomain string `json:"myshopifyDomain"`
		Email           string `json:"email"`
		Plan            struct {
			DisplayName string `json:"displayName"`
		} `json:"plan"`
		BillingAddress struct {
			City          string `json:"city"`
			Country       string `json:"country"`
			CountryCodeV2 string `json:"countryCodeV2"`
		} `json:"billingAddress"`
		CurrencyCode    string `json:"currencyCode"`
		CurrencyFormats struct {
			MoneyFormat string `json:"moneyFormat"`
		} `json:"currencyFormats"`
	} `json:"shop"`
}

type WebhookSubscription struct {
	Id       string `json:"id"`
	Topic    string `json:"topic"`
	Endpoint struct {
		Typename    string `json:"__typename"`
		CallbackUrl string `json:"callbackUrl"`
	} `json:"endpoint"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	ApiVersion struct {
		Handle string `json:"handle"`
	} `json:"apiVersion"`
	Format              string   `json:"format"`
	IncludeFields       []string `json:"includeFields"`
	MetafieldNamespaces []string `json:"metafieldNamespaces"`
}
