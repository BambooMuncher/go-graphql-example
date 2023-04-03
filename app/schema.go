package main

import (
	"sync"
    "github.com/graphql-go/graphql"
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "net/url"
)

/*
 * We are caching the values of this metadata indefinitely. 
 * We could easily set a timer for this or provide a way to refresh it.
 * These values should only change when Blizzard releases a new expansion.
 * For now, we can restart the server to refresh the data.
 */
var cardSetsById map[int]CardSet = nil
var raritiesById map[int]Rarity = nil
var classesById map[int]Class = nil
var cardTypesById map[int]CardType = nil
var minionTypesById map[int]MinionType = nil
var spellSchoolsById map[int]SpellSchool = nil


var rootQuery = graphql.NewObject(graphql.ObjectConfig{
    Name: "Query",
    Fields: graphql.Fields{
        "cards": &graphql.Field{
            Type: graphql.NewList(cardType),
            Args: graphql.FieldConfigArgument{
                "search": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                root := p.Info.RootValue.(map[string]interface{})
                token := root["token"].(string)
                if token == "" {
                    return nil, fmt.Errorf("Bearer token is required.")
                }
                client := &http.Client{}
                result := make(chan CardSearchResponse, 1)
                errors := make(chan error, 100)
                var wg sync.WaitGroup
                wg.Add(7)
                
                go func() {
                    fmt.Printf("Started card search\n")

                    cardSearchResponse, err := performCardSearch(token, *client, p.Args, errors)
                    if err != nil {
                        fmt.Printf("Error performing card search: %s\n", err)
                        errors <- err
                    } 
                    if(cardSearchResponse != nil) {
                        result <- *cardSearchResponse
                    }       

                    fmt.Printf("Completed card search\n")
                    wg.Done()
                }()

                retrieveMetadata(token, *client, &wg, errors)

                fmt.Println("Waiting")

                wg.Wait()

                fmt.Println("Done waiting")

                if len(errors) > 0 {
                    return nil, <- errors
                }

                cardSearchResponse := <-result

                results := make([]interface{}, len(cardSearchResponse.Cards))

                for i, card := range cardSearchResponse.Cards {
                    results[i] = mapCardToGraphQL(card)
                }

                return results, nil
            },
        },
    },
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: rootQuery,
})

func performCardSearch(token string, client http.Client, args map[string]interface{}, errors chan error) (*CardSearchResponse, error) {
    params := url.Values{}
    params.Add("locale", "en_US")
    if(args["search"] != nil) {
        params.Add("textFilter", args["search"].(string))
    }
    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/cards?" + params.Encode(), nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error searching cards: %s\n", err)
        return nil, err
    }
    if resp.StatusCode == 401 {
        return nil, fmt.Errorf("Token is expired or invalid.")
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var cards CardSearchResponse
    if err = json.Unmarshal(body, &cards); err != nil {
        fmt.Printf("Error unmarshalling cards: %s\n", err)
    }

    return &cards, nil
}

func mapCardToGraphQL(card Card) map[string]interface{} {
    var classList []Class = nil
    if card.MultiClassIds != nil && len(card.MultiClassIds) > 0 {
        classList = make([]Class, len(card.MultiClassIds))
        for j, class := range card.MultiClassIds {
            classList[j] = classesById[class]
        }
    } else {
        classList = []Class{classesById[card.ClassId]}
    }

    var minionType *MinionType = nil
    if card.MinionTypeId != nil {
        if value, ok := minionTypesById[*card.MinionTypeId]; ok {
            minionType = &value
        }
    }

    var spellSchool *SpellSchool = nil
    if card.SpellSchoolId != nil {
        if value, ok := spellSchoolsById[*card.SpellSchoolId]; ok {
            spellSchool = &value
        }
    }

    return map[string]interface{}{
        "id": card.ID,
        "name": card.Name,
        "cardSet": cardSetsById[card.CardSetId],
        "rarity": raritiesById[card.RarityId],
        "classes": classList,
        "cardType": cardTypesById[card.CardTypeId],
        "minionType": minionType,
        "spellSchool": spellSchool,
    }
}

func retrieveCardSetsById(token string, client http.Client) (map[int]CardSet, error){
    if cardSetsById != nil {
        return cardSetsById, nil
    }

    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/metadata/sets?locale=en_US", nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error getting sets: %s\n", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var cardSets []CardSet 
    if err = json.Unmarshal(body, &cardSets); err != nil {
        fmt.Printf("Error unmarshalling sets: %s\n", err)
        return nil, err
    }
    cardSetsById = make(map[int]CardSet)
    for _, cardSet := range cardSets {
        cardSetsById[cardSet.ID] = cardSet
    }

    return cardSetsById, nil
}

func retrieveRaritiesById(token string, client http.Client) (map[int]Rarity, error){
    if raritiesById != nil {
        return raritiesById, nil
    }

    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/metadata/rarities?locale=en_US", nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error getting rarities: %s\n", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var rarities []Rarity 
    if err = json.Unmarshal(body, &rarities); err != nil {
        fmt.Printf("Error unmarshalling rarities: %s\n", err)
        return nil, err
    }
    raritiesById = make(map[int]Rarity)
    for _, rarity := range rarities {
        raritiesById[rarity.ID] = rarity
    }

    return raritiesById, nil
}

func retrieveClassesById(token string, client http.Client) (map[int]Class, error){
    if classesById != nil {
        return classesById, nil
    }

    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/metadata/classes?locale=en_US", nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error getting classes: %s\n", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var classes []Class 
    if err = json.Unmarshal(body, &classes); err != nil {
        fmt.Printf("Error unmarshalling classes: %s\n", err)
    }
    classesById = make(map[int]Class)
    for _, class := range classes {
        classesById[class.ID] = class
    }

    return classesById, nil
}

func retrieveCardTypesById(token string, client http.Client) (map[int]CardType, error){
    if cardTypesById != nil {
        return cardTypesById, nil
    }

    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/metadata/types?locale=en_US", nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error getting card types: %s\n", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var cardTypes []CardType 
    if err = json.Unmarshal(body, &cardTypes); err != nil {
        fmt.Printf("Error unmarshalling card types: %s\n", err)
        return nil, err
    }
    cardTypesById = make(map[int]CardType)
    for _, cardType := range cardTypes {
        cardTypesById[cardType.ID] = cardType
    }

    return cardTypesById, nil
}

func retrieveMinionTypesById(token string, client http.Client) (map[int]MinionType, error){
    if minionTypesById != nil {
        return minionTypesById, nil
    }

    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/metadata/minionTypes?locale=en_US", nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error getting minion types: %s\n", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var minionTypes []MinionType 
    if err = json.Unmarshal(body, &minionTypes); err != nil {
        fmt.Printf("Error unmarshalling minion types: %s\n", err)
        return nil, err
    }
    minionTypesById = make(map[int]MinionType)
    for _, minionType := range minionTypes {
        minionTypesById[minionType.ID] = minionType
    }

    return minionTypesById, nil
}

func retrieveSpellSchoolsById(token string, client http.Client) (map[int]SpellSchool, error){
    if spellSchoolsById != nil {
        return spellSchoolsById, nil
    }

    req, err := http.NewRequest("GET", "https://us.api.blizzard.com/hearthstone/metadata/spellSchools?locale=en_US", nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error getting spell schools: %s\n", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var spellSchools []SpellSchool 
    if err = json.Unmarshal(body, &spellSchools); err != nil {
        fmt.Printf("Error unmarshalling spell schools: %s\n", err)
        return nil, err
    }

    spellSchoolsById = make(map[int]SpellSchool)
    for _, spellSchool := range spellSchools {
        spellSchoolsById[spellSchool.ID] = spellSchool
    }

    return spellSchoolsById, nil
}

func retrieveMetadata(token string, client http.Client, wg *sync.WaitGroup, errors chan<- error) {

    go func() {
        _, err := retrieveCardSetsById(token, client)
        if err != nil {
            fmt.Printf("Error retrieving card sets: %s\n", err)
            errors <- err
        } else {
            fmt.Printf("Retrieved card sets\n")
        }
        wg.Done()
    }()

    go func() {
        _, err := retrieveRaritiesById(token, client)
        if err != nil {
            fmt.Printf("Error retrieving rarities: %s\n", err)
            errors <- err
        } else {
            fmt.Printf("Retrieved rarities\n")
        }
        wg.Done()
    }()

    go func() {
        _, err := retrieveClassesById(token, client)
        if err != nil {
            fmt.Printf("Error retrieving classes: %s\n", err)
            errors <- err
        } else {
            fmt.Printf("Retrieved classes\n")
        }
        wg.Done()
    }()

    go func() {
        _, err := retrieveCardTypesById(token, client)
        if err != nil {
            fmt.Printf("Error retrieving card types: %s\n", err)
            errors <- err
        } else {
            fmt.Printf("Retrieved card types\n")
        }
        wg.Done()
    }()

    go func() {
        _, err := retrieveMinionTypesById(token, client)
        if err != nil {
            fmt.Printf("Error retrieving minion types: %s\n", err)
            errors <- err
        } else {
            fmt.Printf("Retrieved minion types\n")
        }
        wg.Done()
    }()

    go func() {
        _, err := retrieveSpellSchoolsById(token, client)
        if err != nil {
            fmt.Printf("Error retrieving spell schools: %s\n", err)
            errors <- err
        } else {
            fmt.Printf("Retrieved spell schools\n")
        }
        wg.Done()
    }()
}
