package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Property représente une propriété (logement) dans la base de données
type Property struct {
	ID                   primitive.ObjectID    `json:"_id" bson:"_id,omitempty"`
	HostID               primitive.ObjectID    `json:"hostId" bson:"hostId" binding:"required"`
	Status               int                   `json:"status" bson:"status"` // 1 = brouillon, 2 = publié
	Slug                 string                `json:"slug" bson:"slug" binding:"required"`
	Name                 string                `json:"name" bson:"name" binding:"required"`
	Description          string                `json:"description,omitempty" bson:"description,omitempty"`
	Address              string                `json:"address" bson:"address" binding:"required"`
	City                 string                `json:"city" bson:"city" binding:"required"`
	Country              string                `json:"country" bson:"country" binding:"required"`
	ZipCode              string                `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
	Images               []string              `json:"images,omitempty" bson:"images,omitempty"`
	CheckInOut           *CheckInOut            `json:"checkInOut" bson:"checkInOut" binding:"required"`
	Wifi                 *Wifi                 `json:"wifi" bson:"wifi" binding:"required"`
	Equipment            *Equipment            `json:"equipment" bson:"equipment" binding:"required"`
	Instructions         *Instructions         `json:"instructions" bson:"instructions" binding:"required"`
	Rules                *Rules                `json:"rules" bson:"rules" binding:"required"`
	Contacts             *Contacts             `json:"contacts" bson:"contacts" binding:"required"`
	LocalRecommendations *LocalRecommendations `json:"localRecommendations" bson:"localRecommendations" binding:"required"`
	Parking              *Parking              `json:"parking" bson:"parking" binding:"required"`
	Transport            *Transport            `json:"transport" bson:"transport" binding:"required"`
	Security             *Security             `json:"security" bson:"security" binding:"required"`
	Services             *Services             `json:"services" bson:"services" binding:"required"`
	BabyKids             *BabyKids            `json:"babyKids" bson:"babyKids" binding:"required"`
	Pets                 *Pets                 `json:"pets" bson:"pets" binding:"required"`
	Entertainment        *Entertainment        `json:"entertainment" bson:"entertainment" binding:"required"`
	Outdoor              *Outdoor             `json:"outdoor" bson:"outdoor" binding:"required"`
	Neighborhood         *Neighborhood        `json:"neighborhood" bson:"neighborhood" binding:"required"`
	Emergency            *Emergency            `json:"emergency" bson:"emergency" binding:"required"`
	CreatedAt            time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time             `json:"updatedAt" bson:"updatedAt"`
	PublishedAt          *time.Time            `json:"publishedAt,omitempty" bson:"publishedAt,omitempty"`
}

// CheckInOut représente les informations d'arrivée et de départ
type CheckInOut struct {
	Enabled             bool   `json:"enabled" bson:"enabled" binding:"required"`
	CheckInTime         string `json:"checkInTime" bson:"checkInTime" binding:"required"` // Format "HH:mm"
	CheckOutTime        string `json:"checkOutTime" bson:"checkOutTime" binding:"required"` // Format "HH:mm"
	SelfCheckIn         bool   `json:"selfCheckIn,omitempty" bson:"selfCheckIn,omitempty"`
	EarlyCheckIn        bool   `json:"earlyCheckIn,omitempty" bson:"earlyCheckIn,omitempty"`
	LateCheckOut        bool   `json:"lateCheckOut,omitempty" bson:"lateCheckOut,omitempty"`
	CheckInInstructions string `json:"checkInInstructions,omitempty" bson:"checkInInstructions,omitempty"`
	CheckOutInstructions string `json:"checkOutInstructions,omitempty" bson:"checkOutInstructions,omitempty"`
	KeyLocation         string `json:"keyLocation,omitempty" bson:"keyLocation,omitempty"`
	AccessCode          string `json:"accessCode,omitempty" bson:"accessCode,omitempty"`
	LockboxCode         string `json:"lockboxCode,omitempty" bson:"lockboxCode,omitempty"`
	BuildingCode        string `json:"buildingCode,omitempty" bson:"buildingCode,omitempty"`
	IntercomCode        string `json:"intercomCode,omitempty" bson:"intercomCode,omitempty"`
	ParkingCode         string `json:"parkingCode,omitempty" bson:"parkingCode,omitempty"`
	GateCode            string `json:"gateCode,omitempty" bson:"gateCode,omitempty"`
}

// Wifi représente les informations Wi-Fi
type Wifi struct {
	Enabled           bool   `json:"enabled" bson:"enabled" binding:"required"`
	NetworkName       string `json:"networkName,omitempty" bson:"networkName,omitempty"`
	Password          string `json:"password,omitempty" bson:"password,omitempty"`
	RouterLocation    string `json:"routerLocation,omitempty" bson:"routerLocation,omitempty"`
	ResetInstructions string `json:"resetInstructions,omitempty" bson:"resetInstructions,omitempty"`
	Notes             string `json:"notes,omitempty" bson:"notes,omitempty"`
}

// EquipmentItem représente un équipement
type EquipmentItem struct {
	ID       string `json:"id" bson:"id" binding:"required"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Category string `json:"category" bson:"category" binding:"required"` // bedroom, bathroom, kitchen, living, outdoor, baby, work, other
}

// Equipment représente les équipements
type Equipment struct {
	Enabled bool            `json:"enabled" bson:"enabled" binding:"required"`
	Items   []EquipmentItem `json:"items" bson:"items" binding:"required"`
}

// InstructionItem représente une consigne spécifique
type InstructionItem struct {
	Enabled bool   `json:"enabled" bson:"enabled" binding:"required"`
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}

// Instructions représente les consignes d'utilisation
type Instructions struct {
	Enabled          bool            `json:"enabled" bson:"enabled" binding:"required"`
	Trash            *InstructionItem `json:"trash,omitempty" bson:"trash,omitempty"`
	Heating          *InstructionItem `json:"heating,omitempty" bson:"heating,omitempty"`
	AirConditioning  *InstructionItem `json:"airConditioning,omitempty" bson:"airConditioning,omitempty"`
	HotWater         *InstructionItem `json:"hotWater,omitempty" bson:"hotWater,omitempty"`
	Appliances       *InstructionItem `json:"appliances,omitempty" bson:"appliances,omitempty"`
	Laundry          *InstructionItem `json:"laundry,omitempty" bson:"laundry,omitempty"`
	Dishwasher       *InstructionItem `json:"dishwasher,omitempty" bson:"dishwasher,omitempty"`
	Oven             *InstructionItem `json:"oven,omitempty" bson:"oven,omitempty"`
	CoffeeMachine    *InstructionItem `json:"coffeeMachine,omitempty" bson:"coffeeMachine,omitempty"`
	Television       *InstructionItem `json:"television,omitempty" bson:"television,omitempty"`
	Sound            *InstructionItem `json:"sound,omitempty" bson:"sound,omitempty"`
	Blinds           *InstructionItem `json:"blinds,omitempty" bson:"blinds,omitempty"`
	Alarm            *InstructionItem `json:"alarm,omitempty" bson:"alarm,omitempty"`
	Safe             *InstructionItem `json:"safe,omitempty" bson:"safe,omitempty"`
	Pool             *InstructionItem `json:"pool,omitempty" bson:"pool,omitempty"`
	Spa              *InstructionItem `json:"spa,omitempty" bson:"spa,omitempty"`
	Garden           *InstructionItem `json:"garden,omitempty" bson:"garden,omitempty"`
	Barbecue         *InstructionItem `json:"barbecue,omitempty" bson:"barbecue,omitempty"`
	Fireplace        *InstructionItem `json:"fireplace,omitempty" bson:"fireplace,omitempty"`
	Other            *InstructionItem `json:"other,omitempty" bson:"other,omitempty"`
}

// Rules représente le règlement intérieur
type Rules struct {
	Enabled         bool     `json:"enabled" bson:"enabled" binding:"required"`
	SmokingAllowed  bool     `json:"smokingAllowed" bson:"smokingAllowed" binding:"required"`
	PetsAllowed     bool     `json:"petsAllowed" bson:"petsAllowed" binding:"required"`
	PartiesAllowed  bool     `json:"partiesAllowed" bson:"partiesAllowed" binding:"required"`
	ChildrenAllowed bool     `json:"childrenAllowed" bson:"childrenAllowed" binding:"required"`
	MaxGuests       *int     `json:"maxGuests,omitempty" bson:"maxGuests,omitempty"`
	QuietHours      string   `json:"quietHours,omitempty" bson:"quietHours,omitempty"`
	HouseRules      []string `json:"houseRules,omitempty" bson:"houseRules,omitempty"`
	AdditionalRules string   `json:"additionalRules,omitempty" bson:"additionalRules,omitempty"`
}

// Contact représente un contact
type Contact struct {
	ID    string `json:"id" bson:"id" binding:"required"`
	Type  string `json:"type" bson:"type" binding:"required"` // host, concierge, cleaning, maintenance, emergency, neighbor, other
	Name  string `json:"name" bson:"name" binding:"required"`
	Phone string `json:"phone" bson:"phone" binding:"required"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Notes string `json:"notes,omitempty" bson:"notes,omitempty"`
}

// Contacts représente les contacts utiles
type Contacts struct {
	Enabled  bool      `json:"enabled" bson:"enabled" binding:"required"`
	Contacts []Contact `json:"contacts" bson:"contacts" binding:"required"`
}

// Recommendation représente une recommandation locale
type Recommendation struct {
	ID          string  `json:"id" bson:"id" binding:"required"`
	Category    string  `json:"category" bson:"category" binding:"required"` // restaurant, cafe, bar, bakery, grocery, market, pharmacy, doctor, hospital, attraction, beach, park, sport, shopping, nightlife, culture, other
	Name        string  `json:"name" bson:"name" binding:"required"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Address     string  `json:"address,omitempty" bson:"address,omitempty"`
	Phone       string  `json:"phone,omitempty" bson:"phone,omitempty"`
	Website     string  `json:"website,omitempty" bson:"website,omitempty"`
	Distance    string  `json:"distance,omitempty" bson:"distance,omitempty"`
	Rating      *float64 `json:"rating,omitempty" bson:"rating,omitempty"` // 1-5
}

// LocalRecommendations représente les recommandations locales
type LocalRecommendations struct {
	Enabled        bool            `json:"enabled" bson:"enabled" binding:"required"`
	Recommendations []Recommendation `json:"recommendations" bson:"recommendations" binding:"required"`
}

// Parking représente les informations de parking
type Parking struct {
	Enabled     bool   `json:"enabled" bson:"enabled" binding:"required"`
	Available   bool   `json:"available,omitempty" bson:"available,omitempty"`
	Type        string `json:"type,omitempty" bson:"type,omitempty"` // street, garage, driveway, private, public
	Free        bool   `json:"free,omitempty" bson:"free,omitempty"`
	Price       string `json:"price,omitempty" bson:"price,omitempty"`
	Instructions string `json:"instructions,omitempty" bson:"instructions,omitempty"`
	AccessCode  string `json:"accessCode,omitempty" bson:"accessCode,omitempty"`
}

// Transport représente les informations de transport
type Transport struct {
	Enabled        bool   `json:"enabled" bson:"enabled" binding:"required"`
	NearestBus     string `json:"nearestBus,omitempty" bson:"nearestBus,omitempty"`
	NearestMetro   string `json:"nearestMetro,omitempty" bson:"nearestMetro,omitempty"`
	NearestTrain   string `json:"nearestTrain,omitempty" bson:"nearestTrain,omitempty"`
	NearestTram    string `json:"nearestTram,omitempty" bson:"nearestTram,omitempty"`
	TaxiInfo       string `json:"taxiInfo,omitempty" bson:"taxiInfo,omitempty"`
	BikeRental     string `json:"bikeRental,omitempty" bson:"bikeRental,omitempty"`
	CarRental      string `json:"carRental,omitempty" bson:"carRental,omitempty"`
	AirportShuttle string `json:"airportShuttle,omitempty" bson:"airportShuttle,omitempty"`
	WalkingInfo    string `json:"walkingInfo,omitempty" bson:"walkingInfo,omitempty"`
}

// Security représente les informations de sécurité
type Security struct {
	Enabled                    bool   `json:"enabled" bson:"enabled" binding:"required"`
	HasAlarm                   bool   `json:"hasAlarm,omitempty" bson:"hasAlarm,omitempty"`
	AlarmCode                  string `json:"alarmCode,omitempty" bson:"alarmCode,omitempty"`
	AlarmInstructions          string `json:"alarmInstructions,omitempty" bson:"alarmInstructions,omitempty"`
	HasSafe                    bool   `json:"hasSafe,omitempty" bson:"hasSafe,omitempty"`
	SafeCode                   string `json:"safeCode,omitempty" bson:"safeCode,omitempty"`
	SafeLocation               string `json:"safeLocation,omitempty" bson:"safeLocation,omitempty"`
	HasFireExtinguisher        bool   `json:"hasFireExtinguisher,omitempty" bson:"hasFireExtinguisher,omitempty"`
	FireExtinguisherLocation   string `json:"fireExtinguisherLocation,omitempty" bson:"fireExtinguisherLocation,omitempty"`
	HasFirstAidKit             bool   `json:"hasFirstAidKit,omitempty" bson:"hasFirstAidKit,omitempty"`
	FirstAidKitLocation        string `json:"firstAidKitLocation,omitempty" bson:"firstAidKitLocation,omitempty"`
	HasSmokeDetector           bool   `json:"hasSmokeDetector,omitempty" bson:"hasSmokeDetector,omitempty"`
	HasCarbonMonoxideDetector  bool   `json:"hasCarbonMonoxideDetector,omitempty" bson:"hasCarbonMonoxideDetector,omitempty"`
	SecurityNotes              string `json:"securityNotes,omitempty" bson:"securityNotes,omitempty"`
}

// Services représente les services proposés
type Services struct {
	Enabled          bool   `json:"enabled" bson:"enabled" binding:"required"`
	LinensIncluded   bool   `json:"linensIncluded,omitempty" bson:"linensIncluded,omitempty"`
	TowelsIncluded   bool   `json:"towelsIncluded,omitempty" bson:"towelsIncluded,omitempty"`
	ToiletryIncluded bool   `json:"toiletryIncluded,omitempty" bson:"toiletryIncluded,omitempty"`
	CleaningIncluded bool   `json:"cleaningIncluded,omitempty" bson:"cleaningIncluded,omitempty"`
	CleaningFrequency string `json:"cleaningFrequency,omitempty" bson:"cleaningFrequency,omitempty"`
	BreakfastIncluded bool   `json:"breakfastIncluded,omitempty" bson:"breakfastIncluded,omitempty"`
	BreakfastDetails  string `json:"breakfastDetails,omitempty" bson:"breakfastDetails,omitempty"`
	ConciergeService  string `json:"conciergeService,omitempty" bson:"conciergeService,omitempty"`
	GroceryDelivery   string `json:"groceryDelivery,omitempty" bson:"groceryDelivery,omitempty"`
	LuggageStorage    bool   `json:"luggageStorage,omitempty" bson:"luggageStorage,omitempty"`
	LaundryService    bool   `json:"laundryService,omitempty" bson:"laundryService,omitempty"`
}

// BabyKids représente les informations bébé & enfants
type BabyKids struct {
	Enabled            bool   `json:"enabled" bson:"enabled" binding:"required"`
	HasCrib            bool   `json:"hasCrib,omitempty" bson:"hasCrib,omitempty"`
	HasHighChair       bool   `json:"hasHighChair,omitempty" bson:"hasHighChair,omitempty"`
	HasBabyGate        bool   `json:"hasBabyGate,omitempty" bson:"hasBabyGate,omitempty"`
	HasChildProofing   bool   `json:"hasChildProofing,omitempty" bson:"hasChildProofing,omitempty"`
	KidsToysAvailable  bool   `json:"kidsToysAvailable,omitempty" bson:"kidsToysAvailable,omitempty"`
	NearbyPlaygrounds  string `json:"nearbyPlaygrounds,omitempty" bson:"nearbyPlaygrounds,omitempty"`
	BabysitterContact  string `json:"babysitterContact,omitempty" bson:"babysitterContact,omitempty"`
	AdditionalInfo     string `json:"additionalInfo,omitempty" bson:"additionalInfo,omitempty"`
}

// Pets représente les informations sur les animaux
type Pets struct {
	Enabled            bool   `json:"enabled" bson:"enabled" binding:"required"`
	PetsAllowed        bool   `json:"petsAllowed,omitempty" bson:"petsAllowed,omitempty"`
	PetFee             string `json:"petFee,omitempty" bson:"petFee,omitempty"`
	PetRules           string `json:"petRules,omitempty" bson:"petRules,omitempty"`
	DogWalkingAreas    string `json:"dogWalkingAreas,omitempty" bson:"dogWalkingAreas,omitempty"`
	NearbyVet          string `json:"nearbyVet,omitempty" bson:"nearbyVet,omitempty"`
	NearbyPetStore     string `json:"nearbyPetStore,omitempty" bson:"nearbyPetStore,omitempty"`
	PetEquipmentAvailable string `json:"petEquipmentAvailable,omitempty" bson:"petEquipmentAvailable,omitempty"`
}

// Entertainment représente les informations de divertissement
type Entertainment struct {
	Enabled            bool   `json:"enabled" bson:"enabled" binding:"required"`
	HasTv              bool   `json:"hasTv,omitempty" bson:"hasTv,omitempty"`
	TvChannels        string `json:"tvChannels,omitempty" bson:"tvChannels,omitempty"`
	HasNetflix        bool   `json:"hasNetflix,omitempty" bson:"hasNetflix,omitempty"`
	NetflixInstructions string `json:"netflixInstructions,omitempty" bson:"netflixInstructions,omitempty"`
	HasSpotify        bool   `json:"hasSpotify,omitempty" bson:"hasSpotify,omitempty"`
	SpotifyInstructions string `json:"spotifyInstructions,omitempty" bson:"spotifyInstructions,omitempty"`
	HasGameConsole    bool   `json:"hasGameConsole,omitempty" bson:"hasGameConsole,omitempty"`
	GameConsoleDetails string `json:"gameConsoleDetails,omitempty" bson:"gameConsoleDetails,omitempty"`
	BoardGames        string `json:"boardGames,omitempty" bson:"boardGames,omitempty"`
	Books             string `json:"books,omitempty" bson:"books,omitempty"`
}

// Outdoor représente les espaces extérieurs
type Outdoor struct {
	Enabled      bool   `json:"enabled" bson:"enabled" binding:"required"`
	HasGarden    bool   `json:"hasGarden,omitempty" bson:"hasGarden,omitempty"`
	GardenInfo   string `json:"gardenInfo,omitempty" bson:"gardenInfo,omitempty"`
	HasTerrace   bool   `json:"hasTerrace,omitempty" bson:"hasTerrace,omitempty"`
	TerraceInfo  string `json:"terraceInfo,omitempty" bson:"terraceInfo,omitempty"`
	HasBalcony   bool   `json:"hasBalcony,omitempty" bson:"hasBalcony,omitempty"`
	BalconyInfo  string `json:"balconyInfo,omitempty" bson:"balconyInfo,omitempty"`
	HasPool      bool   `json:"hasPool,omitempty" bson:"hasPool,omitempty"`
	PoolInfo     string `json:"poolInfo,omitempty" bson:"poolInfo,omitempty"`
	PoolRules    string `json:"poolRules,omitempty" bson:"poolRules,omitempty"`
	HasSpa       bool   `json:"hasSpa,omitempty" bson:"hasSpa,omitempty"`
	SpaInfo      string `json:"spaInfo,omitempty" bson:"spaInfo,omitempty"`
	HasBarbecue  bool   `json:"hasBarbecue,omitempty" bson:"hasBarbecue,omitempty"`
	BarbecueInfo string `json:"barbecueInfo,omitempty" bson:"barbecueInfo,omitempty"`
}

// Neighborhood représente les informations sur le quartier
type Neighborhood struct {
	Enabled          bool   `json:"enabled" bson:"enabled" binding:"required"`
	Description      string `json:"description,omitempty" bson:"description,omitempty"`
	NoiseLevel       string `json:"noiseLevel,omitempty" bson:"noiseLevel,omitempty"` // quiet, moderate, lively
	NeighborInfo     string `json:"neighborInfo,omitempty" bson:"neighborInfo,omitempty"`
	NearbyAttractions string `json:"nearbyAttractions,omitempty" bson:"nearbyAttractions,omitempty"`
	SafetyTips       string `json:"safetyTips,omitempty" bson:"safetyTips,omitempty"`
}

// Emergency représente les informations d'urgence
type Emergency struct {
	Enabled                 bool   `json:"enabled" bson:"enabled" binding:"required"`
	EmergencyNumber         string `json:"emergencyNumber,omitempty" bson:"emergencyNumber,omitempty"` // Défaut: 112
	PoliceNumber            string `json:"policeNumber,omitempty" bson:"policeNumber,omitempty"` // Défaut: 17
	FireNumber              string `json:"fireNumber,omitempty" bson:"fireNumber,omitempty"` // Défaut: 18
	AmbulanceNumber         string `json:"ambulanceNumber,omitempty" bson:"ambulanceNumber,omitempty"` // Défaut: 15
	NearestHospital         string `json:"nearestHospital,omitempty" bson:"nearestHospital,omitempty"`
	NearestHospitalAddress  string `json:"nearestHospitalAddress,omitempty" bson:"nearestHospitalAddress,omitempty"`
	NearestPharmacy         string `json:"nearestPharmacy,omitempty" bson:"nearestPharmacy,omitempty"`
	NearestPharmacyHours    string `json:"nearestPharmacyHours,omitempty" bson:"nearestPharmacyHours,omitempty"`
	DoctorOnCall            string `json:"doctorOnCall,omitempty" bson:"doctorOnCall,omitempty"`
	AdditionalEmergencyInfo string `json:"additionalEmergencyInfo,omitempty" bson:"additionalEmergencyInfo,omitempty"`
}

// CreatePropertyRequest représente la requête pour créer une propriété
type CreatePropertyRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description,omitempty"`
	Address     string   `json:"address" binding:"required"`
	City        string   `json:"city" binding:"required"`
	Country     string   `json:"country" binding:"required"`
	ZipCode     string   `json:"zipCode,omitempty"`
	Images      []string `json:"images,omitempty"`
	// Tous les sous-documents optionnels pour la création
	CheckInOut           *CheckInOut            `json:"checkInOut,omitempty"`
	Wifi                 *Wifi                 `json:"wifi,omitempty"`
	Equipment            *Equipment            `json:"equipment,omitempty"`
	Instructions         *Instructions         `json:"instructions,omitempty"`
	Rules                *Rules                `json:"rules,omitempty"`
	Contacts             *Contacts             `json:"contacts,omitempty"`
	LocalRecommendations *LocalRecommendations `json:"localRecommendations,omitempty"`
	Parking              *Parking              `json:"parking,omitempty"`
	Transport            *Transport            `json:"transport,omitempty"`
	Security             *Security             `json:"security,omitempty"`
	Services             *Services             `json:"services,omitempty"`
	BabyKids             *BabyKids            `json:"babyKids,omitempty"`
	Pets                 *Pets                 `json:"pets,omitempty"`
	Entertainment        *Entertainment        `json:"entertainment,omitempty"`
	Outdoor              *Outdoor             `json:"outdoor,omitempty"`
	Neighborhood         *Neighborhood        `json:"neighborhood,omitempty"`
	Emergency            *Emergency            `json:"emergency,omitempty"`
}

// UpdatePropertyRequest représente la requête pour mettre à jour une propriété (tous les champs optionnels)
type UpdatePropertyRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Address     string   `json:"address,omitempty"`
	City        string   `json:"city,omitempty"`
	Country     string   `json:"country,omitempty"`
	ZipCode     string   `json:"zipCode,omitempty"`
	Images      []string `json:"images,omitempty"`
	// Tous les sous-documents optionnels
	CheckInOut           *CheckInOut            `json:"checkInOut,omitempty"`
	Wifi                 *Wifi                 `json:"wifi,omitempty"`
	Equipment            *Equipment            `json:"equipment,omitempty"`
	Instructions         *Instructions         `json:"instructions,omitempty"`
	Rules                *Rules                `json:"rules,omitempty"`
	Contacts             *Contacts             `json:"contacts,omitempty"`
	LocalRecommendations *LocalRecommendations `json:"localRecommendations,omitempty"`
	Parking              *Parking              `json:"parking,omitempty"`
	Transport            *Transport            `json:"transport,omitempty"`
	Security             *Security             `json:"security,omitempty"`
	Services             *Services             `json:"services,omitempty"`
	BabyKids             *BabyKids            `json:"babyKids,omitempty"`
	Pets                 *Pets                 `json:"pets,omitempty"`
	Entertainment        *Entertainment        `json:"entertainment,omitempty"`
	Outdoor              *Outdoor             `json:"outdoor,omitempty"`
	Neighborhood         *Neighborhood        `json:"neighborhood,omitempty"`
	Emergency            *Emergency            `json:"emergency,omitempty"`
}
