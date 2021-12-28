package model

type Movie struct {
	// Ean,
	Ean string
	// Myyntihinta,
	SalePrice float32
	// Kuntoluokitus,
	Conditions []string
	// Nimi,
	LocalName string
	// Alkuperäinen nimi,
	OriginalName string
	// Vuosi,
	Year uint16
	// Ohjaaja,
	Directors []string
	// Näyttelijöitä,
	Actors []string
	// Formaatti,
	Format string
	// Levyjä,
	Discs uint8
	// Julkaisualue,
	PublicationArea string
	// Julkaisu,
	Publication string
	// Tekstitys,
	Subtitles []string
	// Puhuttu kieli,
	Languages []string
	// Muuta,
	Other string
	// Vuokrajulkaisu,
	IsRental bool
	// Slipcover,
	HasSlipCover bool
	// Kaksipuolinen levy,
	IsTwoSidedDisc bool
	// Skannattu,
	IsReadTested bool
	// Kotelo,
	CaseType string
	// Toimitusluokka
	DeliveryClass string
}
