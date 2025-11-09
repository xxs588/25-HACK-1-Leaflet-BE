package consts

// 头像结构体
type PresetAvatar struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

// 之后换图片URL
var ProfilePictures = []PresetAvatar{
	{ID: 1, URL: "https://example.com/images/profile1.png"},
	{ID: 2, URL: "https://example.com/images/profile2.png"},
	{ID: 3, URL: "https://example.com/images/profile3.png"},
	{ID: 4, URL: "https://example.com/images/profile4.png"},
	{ID: 5, URL: "https://example.com/images/profile5.png"},
}
