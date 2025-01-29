package campaign

import (
	"fmt"
	"strings"
)

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailFormatter struct {
	ID               int                     `json:"id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	ImageURL         string                  `json:"image_url"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	UserID           int                     `json:"user_id"`
	Slug             string                  `json:"slug"`
	Perks            []string                `json:"perks"`
	User             CampaignUserFormatter   `json:"user"`
	Images           []CampaignImageFormater `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormater struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.ImageURL = ""
	campaignFormatter.Slug = campaign.Slug

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	campaignsFormatter := []CampaignFormatter{}

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	CampaignDetailFormatter := CampaignDetailFormatter{}
	CampaignDetailFormatter.ID = campaign.ID
	CampaignDetailFormatter.Name = campaign.Name
	CampaignDetailFormatter.ShortDescription = campaign.ShortDescription
	CampaignDetailFormatter.Description = campaign.Description
	CampaignDetailFormatter.GoalAmount = campaign.GoalAmount
	CampaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	CampaignDetailFormatter.Slug = campaign.Slug
	CampaignDetailFormatter.UserID = campaign.UserID
	CampaignDetailFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		CampaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}
	CampaignDetailFormatter.Perks = perks
	user := campaign.User
	CampaignUserFormatter := CampaignUserFormatter{}
	CampaignUserFormatter.Name = user.Name
	CampaignUserFormatter.ImageURL = user.AvatarFileName

	CampaignDetailFormatter.User = CampaignUserFormatter

	images := []CampaignImageFormater{}

	fmt.Println(campaign.CampaignImages)

	for _, x := range campaign.CampaignImages {
		CampaignImageFormater := CampaignImageFormater{}
		CampaignImageFormater.ImageURL = x.FileName
		isPrimary := false

		if x.IsPrimary == 1 {
			isPrimary = true
		}

		CampaignImageFormater.IsPrimary = isPrimary
		images = append(images, CampaignImageFormater)
	}

	CampaignDetailFormatter.Images = images

	return CampaignDetailFormatter
}
