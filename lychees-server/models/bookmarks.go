package models

type Item struct {
	URL   string `json:"url" binding:"required"`
	Title string `json:"title" bson:"title"`
	Icon  string `json:"icon" bson:"icon"`

	IsSvg bool `json:"isSvg" bson:"isSvg"`

	Turn  bool   `json:"turn"`
	Color string `json:"color,omitempty"`
	//BackgroundColor string `json:"backgroundColor,omitempty" bson:"backgroundColor"`
	IconSize string `json:"iconSize,omitempty" bson:"iconSize,omitempty"`
}

type MoveItem struct {
	Item        Item `json:"item"`
	NewRowIndex uint `json:"newRowIndex"`
	NewColIndex uint `json:"newColIndex"`
}

type IconfontLink struct {
	IconfontLink string `json:"iconfontLink,omitempty" bson:"iconfontLink"`
}

type Bookmarks struct {
	ID uint `json:"-" bson:"_id,omitempty"`
	PersonalInfo
	//更新时间 unix 时间戳
	UpdateAt int64 `json:"updateAt" bson:"updateAt"`

	//WallpaperUrl   string   `json:"wallpaperUrl,omitempty" bson:"wallpaperUrl,omitempty"`
	ArrayBookmarks [][]Item `json:"arrayBookmarks" bson:"arrayBookmarks"`
}
