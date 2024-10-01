package models

type RewildingReferenceLinks struct {
	RewildingReferenceLinksLink          string `bson:"rewilding_reference_links_link,omitempty" json:"rewilding_reference_links_link,omitempty"`
	RewildingReferenceLinksTitle         string `bson:"rewilding_reference_links_title,omitempty" json:"rewilding_reference_links_title,omitempty"`
	RewildingReferenceLinksDescription   string `bson:"rewilding_reference_links_description" json:"rewilding_reference_links_description,omitempty"`
	RewildingReferenceLinksOGTitle       string `bson:"rewilding_reference_links_og_title" json:"rewilding_reference_links_og_title,omitempty"`
	RewildingReferenceLinksOGDescription string `bson:"rewilding_reference_links_og_description" json:"rewilding_reference_links_og_description,omitempty"`
	RewildingReferenceLinksOGImage       string `bson:"rewilding_reference_links_og_image" json:"rewilding_reference_links_og_image,omitempty"`
	RewildingReferenceLinksOGSiteName    string `bson:"rewilding_reference_links_og_site_name" json:"rewilding_reference_links_og_site_name,omitempty"`
}
