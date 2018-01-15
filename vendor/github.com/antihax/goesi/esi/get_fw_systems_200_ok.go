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

/* A list of GetFwSystems200Ok. */
//easyjson:json
type GetFwSystems200OkList []GetFwSystems200Ok

/* 200 ok object */
//easyjson:json
type GetFwSystems200Ok struct {
	SolarSystemId          int32 `json:"solar_system_id,omitempty"`          /* solar_system_id integer */
	OwnerFactionId         int32 `json:"owner_faction_id,omitempty"`         /* owner_faction_id integer */
	OccupierFactionId      int32 `json:"occupier_faction_id,omitempty"`      /* occupier_faction_id integer */
	VictoryPoints          int32 `json:"victory_points,omitempty"`           /* victory_points integer */
	VictoryPointsThreshold int32 `json:"victory_points_threshold,omitempty"` /* victory_points_threshold integer */
	Contested              bool  `json:"contested,omitempty"`                /* contested boolean */
}
