package main

type Geocoder struct {
	Response struct {
		GeoObjectCollection struct {
			MetaDataProperty struct {
				GeocoderResponseMetaData struct {
					Request string `json:"request"`
					Found   string `json:"found"`
					Results string `json:"results"`
				} `json:"GeocoderResponseMetaData"`
			} `json:"metaDataProperty"`
			FeatureMember []*yandexFeatureMember `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

type yandexFeatureMember struct {
	GeoObject struct {
		MetaDataProperty struct {
			GeocoderMetaData struct {
				Kind      string `json:"kind"`
				Text      string `json:"text"`
				Precision string `json:"precision"`
				Address   struct {
					CountryCode string `json:"country_code"`
					PostalCode  string `json:"postal_code"`
					Formatted   string `json:"formatted"`
					Components  []struct {
						Kind string `json:"kind"`
						Name string `json:"name"`
					} `json:"Components"`
				} `json:"Address"`
			} `json:"GeocoderMetaData"`
		} `json:"metaDataProperty"`
		Description string `json:"description"`
		Name        string `json:"name"`
		BoundedBy   struct {
			Envelope struct {
				LowerCorner string `json:"lowerCorner"`
				UpperCorner string `json:"upperCorner"`
			} `json:"Envelope"`
		} `json:"boundedBy"`
		Point struct {
			Pos string `json:"pos"`
		} `json:"Point"`
	} `json:"GeoObject"`
}
