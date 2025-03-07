package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mjande/meals-microservice/database"
	"github.com/mjande/meals-microservice/utils"
)

type Meal struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	RecipeID int64  `json:"recipe_id"`
	Date     string `json:"date"`
}

func ListMealsByDate(ctx context.Context, start, end string) ([]Meal, error) {
	userId := utils.ExtractUserIDFromContext(ctx)
	query := `SELECT id, user_id, recipe_id, date FROM meals where user_id = $1 AND date BETWEEN $2 AND $3`

	rows, err := database.DB.Query(ctx, query, userId, start, end)
	if err != nil && err == pgx.ErrNoRows {
		return []Meal{}, nil
	} else if err != nil {
		return []Meal{}, err
	}

	var meals []Meal
	for rows.Next() {
		var meal Meal

		err = rows.Scan(&meal.ID, &meal.UserID, &meal.RecipeID, &meal.Date)
		if err != nil {
			return []Meal{}, err
		}

		meals = append(meals, meal)
	}

	return meals, nil
}

func FindMeal(ctx context.Context, id int64) (Meal, error) {
	query := `SELECT id, user_id, recipe_id, date FROM meals WHERE id = ?`

	result := database.DB.QueryRow(ctx, query, id)

	var meal Meal
	err := result.Scan(&meal.ID, &meal.UserID, &meal.RecipeID, &meal.Date)
	if err != nil {
		return Meal{}, err
	}

	return meal, nil
}

func CreateMeal(ctx context.Context, meal Meal) (int64, error) {
	userId := utils.ExtractUserIDFromContext(ctx)
	query := `INSERT INTO meals (user_id, recipe_id, date) VALUES ($1, $2, $3) RETURNING id`

	row := database.DB.QueryRow(ctx, query, userId, meal.RecipeID, meal.Date)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func DeleteMeal(ctx context.Context, id int64) error {
	query := `DELETE FROM meals WHERE id = ?`

	_, err := database.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
