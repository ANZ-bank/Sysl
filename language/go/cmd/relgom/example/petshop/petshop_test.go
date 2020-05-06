package petshopmodel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPetShopModel(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		_ = NewPetShopModel()
	})
}

func TestInsert(t *testing.T) {
	t.Parallel()

	a := NewPetShopModel()
	b, lab, err := a.GetBreed().Insert().WithBreedName("Labrador").WithNumLegs(4).Apply()
	require.NoError(t, err)
	_, pet, err := b.GetPet().Insert().WithBreed(lab).WithName("Dingo").WithDob(time.Now()).Apply()
	require.NoError(t, err)
	assert.Equal(t, "Labrador", *lab.BreedName())
	assert.Equal(t, "Dingo", *pet.Name())
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	a := NewPetShopModel()
	b, lab, err := a.GetBreed().Insert().WithBreedName("Labrador").WithNumLegs(4).Apply()
	require.NoError(t, err)
	c, pet, err := b.GetPet().Insert().WithBreed(lab).WithName("Dingo").WithNumLegs(*lab.NumLegs()).Apply()
	require.NoError(t, err)
	assert.EqualValues(t, 4, *pet.NumLegs())
	_, pet, err = c.GetPet().Update(pet).WithNumLegs(3).Apply()
	require.NoError(t, err)
	assert.EqualValues(t, 3, *pet.NumLegs())
}

func TestDelete(t *testing.T) {
	t.Parallel()

	a := NewPetShopModel()
	b, lab, err := a.GetBreed().Insert().WithBreedName("Labrador").WithNumLegs(4).Apply()
	require.NoError(t, err)
	c, dob, err := b.GetBreed().Insert().WithBreedName("Doberman").WithNumLegs(4).Apply()
	require.NoError(t, err)
	require.NotPanics(t, func() { _, err = c.GetBreed().Delete(lab) })
	require.NotPanics(t, func() { _, err = c.GetBreed().Delete(dob) })
}

func TestLookup(t *testing.T) {
	t.Parallel()

	a := NewPetShopModel()
	b, lab, err := a.GetBreed().Insert().WithBreedName("Labrador").WithNumLegs(4).Apply()
	require.NoError(t, err)
	c, dob, err := b.GetBreed().Insert().WithBreedName("Doberman").WithNumLegs(4).Apply()
	require.NoError(t, err)

	_, has := a.GetBreed().Lookup(lab.BreedID())
	assert.False(t, has)
	labLU, has := b.GetBreed().Lookup(lab.BreedID())
	if assert.True(t, has) {
		assert.Equal(t, lab.BreedID(), labLU.BreedID())
	}
	labLU, has = c.GetBreed().Lookup(lab.BreedID())
	if assert.True(t, has) {
		assert.Equal(t, lab.BreedID(), labLU.BreedID())
	}

	_, has = a.GetBreed().Lookup(dob.BreedID())
	assert.False(t, has)
	_, has = b.GetBreed().Lookup(dob.BreedID())
	assert.False(t, has)
	dobLU, has := c.GetBreed().Lookup(dob.BreedID())
	if assert.True(t, has) {
		assert.Equal(t, dob.BreedID(), dobLU.BreedID())
	}
}
