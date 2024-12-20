// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/brands/all": {
            "get": {
                "description": "Возвращает все бренды с возможностью фильтрации и сортировки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Получение всех брендов",
                "responses": {
                    "200": {
                        "description": "Список брендов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Brand"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to fetch brand",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/brands/create": {
            "post": {
                "description": "Эндпоинт для создания нового бренда",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Создание нового бренда",
                "parameters": [
                    {
                        "description": "Данные нового бренда",
                        "name": "brand",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Brand"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Brand created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to create brand",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/brands/delete/{id}": {
            "delete": {
                "description": "Выполняет мягкое удаление бренда по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Мягкое удаление бренда",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID бренда",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Brand soft-deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to delete brand",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/brands/filter": {
            "get": {
                "description": "Возвращает все бренды с возможностью фильтрации и сортировки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Фильтрация брендов",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Фильтр по имени бренда",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по стране происхождения",
                        "name": "origin_country",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Фильтр по популярности (целое число)",
                        "name": "popularity",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Фильтр по признаку премиум-бренда",
                        "name": "is_premium",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Фильтр по признаку предстоящего бренда",
                        "name": "is_upcoming",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Фильтр по году основания",
                        "name": "founded_year",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поле сортировки (например, 'name', '-popularity')",
                        "name": "sort",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список брендов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Brand"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to fetch brands",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/brands/restore/{id}": {
            "post": {
                "description": "Восстановление бренда, который был удалён ранее (мягкое удаление)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Восстановление мягко удалённого бренда",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID бренда",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Brand restored successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to restore brand",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/brands/update/{id}": {
            "put": {
                "description": "Обновление данных бренда по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Обновление бренда по ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID бренда (UUIDv7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Обновлённые данные бренда",
                        "name": "brand",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Brand"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Brand updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Brand not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to update brand",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/brands/{id}": {
            "get": {
                "description": "Получение информации о бренде по его уникальному ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brand"
                ],
                "summary": "Получение бренда по ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID бренда",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Бренд найден",
                        "schema": {
                            "$ref": "#/definitions/dto.Brand"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Brand not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/all": {
            "get": {
                "description": "Возвращает все модели с возможностью фильтрации и сортировки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Получение всех моделей",
                "responses": {
                    "200": {
                        "description": "Список моделей",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Model"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to fetch models",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/create": {
            "post": {
                "description": "Эндпоинт для создания новой модели",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Создание новой модели",
                "parameters": [
                    {
                        "description": "Данные новой модели",
                        "name": "model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Model"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Model created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to create model",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/delete/{id}": {
            "delete": {
                "description": "Выполняет мягкое удаление модели по её ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Мягкое удаление модели",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID модели (UUIDv7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Model soft-deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to delete model",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/filter": {
            "get": {
                "description": "Возвращает все модели с возможностью фильтрации и сортировки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Фильтрация моделей",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Фильтр по имени модели",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по идентификатору бренда",
                        "name": "brand_id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Фильтр по популярности (целое число)",
                        "name": "popularity",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Фильтр по признаку премиум-модели",
                        "name": "is_limited",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поле сортировки (например, 'name', '-popularity')",
                        "name": "sort",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список моделей",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Model"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to fetch models",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/restore/{id}": {
            "post": {
                "description": "Восстановление модели, которая была удалена ранее",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Восстановление мягко удалённой модели",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID модели (UUIDv7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Model restored successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to restore model",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/update/{id}": {
            "put": {
                "description": "Обновление данных модели по её ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Обновление модели по ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID модели (UUIDv7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Обновлённые данные модели",
                        "name": "model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Model"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Model updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Model not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to update model",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/models/{id}": {
            "get": {
                "description": "Получение информации о модели по её уникальному ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "models"
                ],
                "summary": "Получение модели по ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID модели (UUIDv7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Модель найдена",
                        "schema": {
                            "$ref": "#/definitions/dto.Model"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Model not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.Brand": {
            "type": "object",
            "properties": {
                "cover_image_url": {
                    "description": "URL обложки",
                    "type": "string"
                },
                "created_at": {
                    "description": "Время создания",
                    "type": "string"
                },
                "description": {
                    "description": "Описание/история бренда",
                    "type": "string"
                },
                "founded_year": {
                    "description": "Год основания",
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "is_deleted": {
                    "description": "Флаг удаления",
                    "type": "boolean"
                },
                "is_premium": {
                    "description": "Флаг премиального бренда",
                    "type": "boolean"
                },
                "is_upcoming": {
                    "description": "Флаг \"Скоро\"",
                    "type": "boolean"
                },
                "link": {
                    "description": "Бренд на английском",
                    "type": "string"
                },
                "logo_url": {
                    "description": "URL логотипа",
                    "type": "string"
                },
                "name": {
                    "description": "Название бренда",
                    "type": "string"
                },
                "origin_country": {
                    "description": "Страна происхождения",
                    "type": "string"
                },
                "popularity": {
                    "description": "Индекс популярности",
                    "type": "integer"
                },
                "updated_at": {
                    "description": "Время обновления",
                    "type": "string"
                }
            }
        },
        "dto.Model": {
            "type": "object",
            "properties": {
                "brand_id": {
                    "type": "string"
                },
                "created_at": {
                    "description": "Время создания",
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_deleted": {
                    "description": "Флаг удаления",
                    "type": "boolean"
                },
                "is_limited": {
                    "description": "Флаг ограниченного выпуска",
                    "type": "boolean"
                },
                "is_upcoming": {
                    "description": "Флаг \"Скоро\"",
                    "type": "boolean"
                },
                "name": {
                    "description": "Название модели",
                    "type": "string"
                },
                "release_date": {
                    "description": "Дата релиза",
                    "type": "string"
                },
                "updated_at": {
                    "description": "Время обновления",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Brands API",
	Description:      "Микро-сервис для работы с брендами и моделями",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
