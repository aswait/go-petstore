package models

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type StatusForm struct {
	Statuses []string `form:"status"`
}

type TagsForm struct {
	Tags []string `form:"tags"`
}

type PetIdForm struct {
	Name   string `form:"name"`
	Status string `form:"status"`
}

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	UserStatus int    `json:"userStatus"`
}

type Order struct {
	ID       int    `json:"id"`
	PetID    int    `json:"petId"`
	Pet      Pet    `gorm:"foreignKey:PetID"`
	Quantity int    `json:"quantity"`
	ShipDate string `json:"shipDate"`
	Status   string `json:"status"`
	Complete bool   `json:"complete"`
}

type OrderResponse struct {
	ID       int    `json:"id"`
	PetID    int    `json:"petId"`
	Quantity int    `json:"quantity"`
	ShipDate string `json:"shipDate"`
	Status   string `json:"status"`
	Complete bool   `json:"complete"`
}

type Pet struct {
	ID         int `json:"id"`
	CategoryID int
	Category   Category   `json:"category" gorm:"foreignKey:CategoryID"`
	Name       string     `json:"name"`
	PhotoUrls  []PhotoUrl `json:"-" gorm:"foreignKey:PetReferID"`
	Tags       []Tag      `json:"tags" gorm:"many2many:pet_tags;constraint:OnDelete:CASCADE"`
	Status     string     `json:"status"`
}

type PetJSON struct {
	ID        int      `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status"`
}

type PhotoUrl struct {
	ID         int    `json:"-"`
	PhotoUrl   string `json:"photo_url"`
	PetReferID int    `json:"-"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PetsStatuses struct {
	Available int `json:"available"`
	Pending   int `json:"pending"`
	Sold      int `json:"sold"`
}
