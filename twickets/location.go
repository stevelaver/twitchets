package twickets

import (
	"encoding/json"
	"fmt"

	"github.com/orsinium-labs/enum"
)

type Country enum.Member[string]

func (d *Country) UnmarshalJSON(data []byte) error {
	var countryString string
	err := json.Unmarshal(data, &countryString)
	if err != nil {
		return err
	}

	country := Countries.Parse(countryString)
	if country == nil {
		return fmt.Errorf("country '%s' is not valid", countryString)
	}
	*d = *country
	return nil
}

var (
	// Countries
	country = enum.NewBuilder[string, Country]()

	CountryUnitedKingdom = country.Add(Country{"GB"})

	Countries = country.Enum()
)

const countryQueryKey = "countryCode"

type Region enum.Member[string]

func (r *Region) UnmarshalJSON(data []byte) error {
	var regionBytes []byte
	err := json.Unmarshal(data, &regionBytes)
	if err != nil {
		return err
	}
	return r.UnmarshalText(regionBytes)
}

func (r *Region) UnmarshalText(data []byte) error {
	regionString := string(data)
	region := Regions.Parse(regionString)
	if region == nil {
		return fmt.Errorf("region '%s' is not valid", regionString)
	}
	*r = *region
	return nil
}

var (
	region = enum.NewBuilder[string, Region]()

	RegionEastAnglia     = region.Add(Region{"GBEA"})
	RegionLondon         = region.Add(Region{"GBLO"})
	RegionMidlands       = region.Add(Region{"GBMI"})
	RegionNorth          = region.Add(Region{"GBNO"})
	RegionNorthEast      = region.Add(Region{"GBNE"})
	RegionNorthernIsland = region.Add(Region{"GBNI"})
	RegionNorthWest      = region.Add(Region{"GBNW"})
	RegionScotland       = region.Add(Region{"GBSC"})
	RegionSouth          = region.Add(Region{"GBSO"})
	RegionSouthEast      = region.Add(Region{"GBSE"})
	RegionSouthWest      = region.Add(Region{"GBSW"})
	RegionWales          = region.Add(Region{"GBWA"})

	Regions = region.Enum()
)

const regionQueryKey = "regionCode"
