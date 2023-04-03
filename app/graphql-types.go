package main

import (
    "github.com/graphql-go/graphql"
)

var cardType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Card",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "cardSet": &graphql.Field{
            Type: cardSetType,
        },
        "rarity": &graphql.Field{
            Type: rarityType,
        },
        "classes": &graphql.Field{
            Type: graphql.NewList(classType),
        },
        "cardType": &graphql.Field{
            Type: cardTypeType,
        },
        "minionType": &graphql.Field{
            Type: minionType,
        },
        "spellSchool": &graphql.Field{
            Type: spellSchoolType,
        },
    },
})

var cardSetType = graphql.NewObject(graphql.ObjectConfig{
    Name: "CardSet",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "slug": &graphql.Field{
            Type: graphql.String,
        },
        "type": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var rarityType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Rarity",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "slug": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var classType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Class",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "slug": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var cardTypeType = graphql.NewObject(graphql.ObjectConfig{
    Name: "CardType",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "slug": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var minionType = graphql.NewObject(graphql.ObjectConfig{
    Name: "MinionType",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "slug": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var spellSchoolType = graphql.NewObject(graphql.ObjectConfig{
    Name: "SpellSchool",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
        "slug": &graphql.Field{
            Type: graphql.String,
        },
    },
})