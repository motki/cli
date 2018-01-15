/*
 * EVE Swagger Interface
 *
 * An OpenAPI for EVE Online
 *
 * OpenAPI spec version: 0.7.5
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package esi

import (
	"time"
)

/* A list of GetMarketsStructuresStructureId200Ok. */
//easyjson:json
type GetMarketsStructuresStructureId200OkList []GetMarketsStructuresStructureId200Ok

/* 200 ok object */
//easyjson:json
type GetMarketsStructuresStructureId200Ok struct {
	OrderId      int64     `json:"order_id,omitempty"`      /* order_id integer */
	TypeId       int32     `json:"type_id,omitempty"`       /* type_id integer */
	LocationId   int64     `json:"location_id,omitempty"`   /* location_id integer */
	VolumeTotal  int32     `json:"volume_total,omitempty"`  /* volume_total integer */
	VolumeRemain int32     `json:"volume_remain,omitempty"` /* volume_remain integer */
	MinVolume    int32     `json:"min_volume,omitempty"`    /* min_volume integer */
	Price        float64   `json:"price,omitempty"`         /* price number */
	IsBuyOrder   bool      `json:"is_buy_order,omitempty"`  /* is_buy_order boolean */
	Duration     int32     `json:"duration,omitempty"`      /* duration integer */
	Issued       time.Time `json:"issued,omitempty"`        /* issued string */
	Range_       string    `json:"range,omitempty"`         /* range string */
}
