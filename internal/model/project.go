package model

type Project struct {
	ID       uint     `json:"id"     gorm:"primaryKey"`
	Title    string   `json:"title"  gorm:"not null"`
	ClientID uint     `json:"-"      gorm:"not null"`
	PmID     uint     `json:"-"      gorm:"not null"`
	Client   User     `json:"client" gorm:"foreignKey:ClientID;references:ID"`
	Pm       User     `json:"pm"     gorm:"foreignKey:PmID;references:ID"`
	Reports  []Report `json:"-"      gorm:"foreignKey:ProjectID"`
}

type ProjectInput struct {
	Title    string `json:"title"`
	ClientID uint   `json:"client_id"`
	PmID     uint   `json:"pm_id"`
}

type ProjectResponse struct {
	ID     uint         `json:"id"`
	Title  string       `json:"title"`
	Client UserResponse `json:"client"`
	Pm     UserResponse `json:"pm"`
}

type ProjectResponseFull struct {
	ID      uint         `json:"id"`
	Title   string       `json:"title"`
	Client  UserResponse `json:"client"`
	Pm      UserResponse `json:"pm"`
	Reports []Report     `json:"reports"`
}
