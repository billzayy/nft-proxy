package nft_proxy

import "time"

type Media struct {
	ID              uint      `json:"-" gorm:"primaryKey"`
	Mint            string    `json:"mint" gorm:"uniqueIndex"`
	MintDecimals    uint8     `json:"decimals"`
	ImageUri        string    `json:"imageUri"`
	ImageType       string    `json:"imageType"`
	MediaUri        string    `json:"mediaUri,omitempty"`
	MediaType       string    `json:"mediaType,omitempty"`
	LocalPath       string    `json:"-"`
	Name            string    `json:"name,omitempty"`
	Symbol          string    `json:"symbol,omitempty"`
	UpdateAuthority string    `json:"updateAuthority,omitempty"`
	CreatedAt       time.Time `json:"-"`
}

type SolanaMedia struct {
	ID              uint      `json:"-" gorm:"primaryKey"`
	Mint            string    `json:"mint" gorm:"uniqueIndex"`
	MintDecimals    uint8     `json:"decimals"`
	ImageUri        string    `json:"imageUri"`
	ImageType       string    `json:"imageType"`
	MediaUri        string    `json:"mediaUri"`
	MediaType       string    `json:"mediaType"`
	LocalPath       string    `json:"-"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	UpdateAuthority string    `json:"updateAuthority"`
	CreatedAt       time.Time `json:"-"`
}

func (m *SolanaMedia) Media() *Media {
	return &Media{
		ID:              m.ID,
		Mint:            m.Mint,
		MintDecimals:    m.MintDecimals,
		ImageUri:        m.ImageUri,
		ImageType:       m.ImageType,
		MediaUri:        m.MediaUri,
		MediaType:       m.MediaType,
		LocalPath:       m.LocalPath,
		Name:            m.Name,
		Symbol:          m.Symbol,
		UpdateAuthority: m.UpdateAuthority,
		CreatedAt:       m.CreatedAt,
	}
}
