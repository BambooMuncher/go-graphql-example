package main

type Card struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Collectible int `json:"collectible"`
    ClassId int `json:"classId"`
    MultiClassIds []int `json:"multiClassIds"`
    SpellSchoolId *int `json:"spellSchoolId"`
    CardTypeId int `json:"cardTypeId"`
    CardSetId int `json:"cardSetId"`
    RarityId int `json:"rarityId"`
    MinionTypeId *int `json:"minionTypeId"`
    ArtistName string `json:"artistName"`
    ManaCost int `json:"manaCost"`
    Attack int `json:"attack"`
    Health int `json:"health"`
    Text string `json:"text"`
    Image string `json:"image"`
    ImageGold string `json:"imageGold"`
    FlavorText string `json:"flavorText"`
    CropImage string `json:"cropImage"`
    ParentId int `json:"parentId"`
    KeywordIds []int `json:"keywordIds"`
}

type CardSet struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
    Type string `json:"type"`
    CollectibleCount int `json:"collectibleCount"`
    CollectibleRevealedCount int `json:"collectibleRevealedCount"`
    NonCollectibleCount int `json:"nonCollectibleCount"`
    NonCollectibleRevealedCount int `json:"nonCollectibleRevealedCount"`
}

type Rarity struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

type Class struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

type CardType struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

type MinionType struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}
type SpellSchool struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

type CardSearchResponse struct {
    Cards []Card `json:"cards"`
    CardCount int `json:"cardCount"`
    PageCount int `json:"pageCount"`
    Page int `json:"page"`
}
