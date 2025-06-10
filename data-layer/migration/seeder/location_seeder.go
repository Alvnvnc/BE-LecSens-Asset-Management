package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type LocationSeeder struct {
	repo *repository.LocationRepository
}

func NewLocationSeeder(db *sql.DB) *LocationSeeder {
	return &LocationSeeder{
		repo: repository.NewLocationRepository(db),
	}
}

// locationData contains realistic location data for Indonesian regions
var locationData = []entity.Location{
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000001"),
		RegionCode:     "DKI-01",
		Name:           "Jakarta Pusat",
		Description:    "Wilayah pusat pemerintahan DKI Jakarta",
		Address:        "Jakarta Pusat, DKI Jakarta, Indonesia",
		Longitude:      106.8451,
		Latitude:       -6.1751,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000002"),
		RegionCode:     "DKI-02",
		Name:           "Jakarta Selatan",
		Description:    "Wilayah Jakarta Selatan",
		Address:        "Jakarta Selatan, DKI Jakarta, Indonesia",
		Longitude:      106.8294,
		Latitude:       -6.2615,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000003"),
		RegionCode:     "DKI-03",
		Name:           "Jakarta Utara",
		Description:    "Wilayah Jakarta Utara",
		Address:        "Jakarta Utara, DKI Jakarta, Indonesia",
		Longitude:      106.8784,
		Latitude:       -6.1388,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000004"),
		RegionCode:     "JB-01",
		Name:           "Bandung",
		Description:    "Kota Bandung, Jawa Barat",
		Address:        "Bandung, Jawa Barat, Indonesia",
		Longitude:      107.6191,
		Latitude:       -6.9175,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000005"),
		RegionCode:     "JT-01",
		Name:           "Surabaya",
		Description:    "Kota Surabaya, Jawa Timur",
		Address:        "Surabaya, Jawa Timur, Indonesia",
		Longitude:      112.7521,
		Latitude:       -7.2575,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000006"),
		RegionCode:     "YGY-01",
		Name:           "Yogyakarta",
		Description:    "Kota Yogyakarta, DIY",
		Address:        "Yogyakarta, DIY, Indonesia",
		Longitude:      110.3695,
		Latitude:       -7.7956,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000007"),
		RegionCode:     "JT-02",
		Name:           "Malang",
		Description:    "Kota Malang, Jawa Timur",
		Address:        "Malang, Jawa Timur, Indonesia",
		Longitude:      112.6304,
		Latitude:       -7.9666,
		HierarchyLevel: 2,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000008"),
		RegionCode:     "SUM-01",
		Name:           "Medan",
		Description:    "Kota Medan, Sumatera Utara",
		Address:        "Medan, Sumatera Utara, Indonesia",
		Longitude:      98.6748,
		Latitude:       3.5952,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000009"),
		RegionCode:     "BAL-01",
		Name:           "Denpasar",
		Description:    "Kota Denpasar, Bali",
		Address:        "Denpasar, Bali, Indonesia",
		Longitude:      115.2126,
		Latitude:       -8.6705,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("01234567-1111-1111-1111-000000000010"),
		RegionCode:     "KSEL-01",
		Name:           "Makassar",
		Description:    "Kota Makassar, Sulawesi Selatan",
		Address:        "Makassar, Sulawesi Selatan, Indonesia",
		Longitude:      119.4221,
		Latitude:       -5.1477,
		HierarchyLevel: 1,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
}

func (s *LocationSeeder) Seed(ctx context.Context) error {
	log.Println("Starting Location seeding...")

	// Insert locations
	for i, location := range locationData {
		locationCopy := location // Create a copy to avoid pointer issues
		err := s.repo.Create(ctx, &locationCopy)
		if err != nil {
			return fmt.Errorf("failed to create location %d (%s): %w", i+1, location.Name, err)
		}
		log.Printf("Created location: %s (%s)", location.Name, location.RegionCode)
	}

	log.Printf("Successfully seeded %d locations", len(locationData))
	return nil
}

func (s *LocationSeeder) GetLocationIDs() []uuid.UUID {
	var ids []uuid.UUID
	for _, location := range locationData {
		ids = append(ids, location.ID)
	}
	return ids
}
