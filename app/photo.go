package app

type GetPhoto struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

type UploadPhoto struct {
	Title    string `json:"title" valid:"required"`
	Caption  string `json:"caption" valid:"required"`
	PhotoUrl string `json:"photo_url" valid:"required"`
}

type UpdatePhoto struct {
	Title    string `json:"title" valid:"required"`
	Caption  string `json:"caption" valid:"required"`
	PhotoUrl string `json:"photo_url" valid:"required"`
}