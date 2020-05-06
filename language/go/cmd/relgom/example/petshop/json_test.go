package petshopmodel

import (
	"encoding/json"
	"testing"

	"github.com/anz-bank/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSONEmptyModel(t *testing.T) {
	t.Parallel()

	m := NewPetShopModel()
	json, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(json))
}

func TestMarshalJSONSingleTuple(t *testing.T) {
	t.Parallel()

	lifespan := decimal.MustParse64("12.25")

	m := NewPetShopModel()

	m, lab, err := m.GetBreed().Insert().
		WithBreedName("Labrador").
		WithNumLegs(4).
		WithLegRank(0).
		WithSpecies("Dog").
		WithAvgLifespan(lifespan).
		Apply()
	require.NoError(t, err)
	assert.EqualValues(t, "Labrador", *lab.BreedName())
	assert.EqualValues(t, 4, *lab.NumLegs())
	assert.EqualValues(t, 0, *lab.LegRank())
	assert.EqualValues(t, "Dog", *lab.Species())
	assert.EqualValues(t, lifespan, *lab.AvgLifespan())

	json, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, `{"Breed":[{
		"breedId":1,
		"breedName":"Labrador",
		"numLegs":4,
		"legRank": 0,
		"species": "Dog",
		"avgLifespan": 12.25
	}]}`, string(json))
}

func TestMarshalJSONTwoTuples(t *testing.T) {
	t.Parallel()

	m := NewPetShopModel()

	m, _, err := m.GetBreed().Insert().WithBreedName("Labrador").WithNumLegs(4).Apply()
	require.NoError(t, err)
	m, _, err = m.GetBreed().Insert().WithBreedName("Doberman").WithNumLegs(4).Apply()
	require.NoError(t, err)

	j, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, canonicalJSONString(`{"Breed":[
		{"breedId":1,"breedName":"Labrador","numLegs":4},
		{"breedId":2,"breedName":"Doberman","numLegs":4}
	]}`), canonicalJSONBytes(j))
}

func twoBreedsAndTwoPets() string {
	return `{
		"Breed":[
			{"breedId":1,"breedName":"Labrador","numLegs":4},
			{"breedId":2,"breedName":"Doberman","numLegs":4}
		],
		"Pet":[
			{"petId":3,"breedId":1,"name":"Scruffy","numLegs":4,"desexed":false},
			{"petId":4,"breedId":2,"name":"Napoleon","numLegs":3,"desexed":true}
		]
	}`
}

func TestMarshalJSONTwoRelations(t *testing.T) {
	t.Parallel()

	m := NewPetShopModel()

	m, lab, err := m.GetBreed().Insert().WithBreedName("Labrador").WithNumLegs(4).Apply()
	require.NoError(t, err)
	m, dob, err := m.GetBreed().Insert().WithBreedName("Doberman").WithNumLegs(4).Apply()
	require.NoError(t, err)

	m, _, err = m.GetPet().Insert().WithBreed(lab).WithName("Scruffy").WithNumLegs(4).WithDesexed(false).Apply()
	require.NoError(t, err)
	m, _, err = m.GetPet().Insert().WithBreed(dob).WithName("Napoleon").WithNumLegs(3).WithDesexed(true).Apply()
	require.NoError(t, err)

	j, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, canonicalJSONString(twoBreedsAndTwoPets()), canonicalJSONBytes(j))
}

func TestUnmarshalJSONTwoRelations(t *testing.T) {
	t.Parallel()

	m := NewPetShopModel()
	require.NoError(t, json.Unmarshal([]byte(twoBreedsAndTwoPets()), &m))

	j, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, canonicalJSONString(twoBreedsAndTwoPets()), canonicalJSONBytes(j))
}
