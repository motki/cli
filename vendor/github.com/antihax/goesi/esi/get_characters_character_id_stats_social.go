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

/* A list of GetCharactersCharacterIdStatsSocial. */
//easyjson:json
type GetCharactersCharacterIdStatsSocialList []GetCharactersCharacterIdStatsSocial

/* social object */
//easyjson:json
type GetCharactersCharacterIdStatsSocial struct {
	AddContactBad             int64 `json:"add_contact_bad,omitempty"`             /* add_contact_bad integer */
	AddContactGood            int64 `json:"add_contact_good,omitempty"`            /* add_contact_good integer */
	AddContactHigh            int64 `json:"add_contact_high,omitempty"`            /* add_contact_high integer */
	AddContactHorrible        int64 `json:"add_contact_horrible,omitempty"`        /* add_contact_horrible integer */
	AddContactNeutral         int64 `json:"add_contact_neutral,omitempty"`         /* add_contact_neutral integer */
	AddNote                   int64 `json:"add_note,omitempty"`                    /* add_note integer */
	AddedAsContactBad         int64 `json:"added_as_contact_bad,omitempty"`        /* added_as_contact_bad integer */
	AddedAsContactGood        int64 `json:"added_as_contact_good,omitempty"`       /* added_as_contact_good integer */
	AddedAsContactHigh        int64 `json:"added_as_contact_high,omitempty"`       /* added_as_contact_high integer */
	AddedAsContactHorrible    int64 `json:"added_as_contact_horrible,omitempty"`   /* added_as_contact_horrible integer */
	AddedAsContactNeutral     int64 `json:"added_as_contact_neutral,omitempty"`    /* added_as_contact_neutral integer */
	CalendarEventCreated      int64 `json:"calendar_event_created,omitempty"`      /* calendar_event_created integer */
	ChatMessagesAlliance      int64 `json:"chat_messages_alliance,omitempty"`      /* chat_messages_alliance integer */
	ChatMessagesConstellation int64 `json:"chat_messages_constellation,omitempty"` /* chat_messages_constellation integer */
	ChatMessagesCorporation   int64 `json:"chat_messages_corporation,omitempty"`   /* chat_messages_corporation integer */
	ChatMessagesFleet         int64 `json:"chat_messages_fleet,omitempty"`         /* chat_messages_fleet integer */
	ChatMessagesRegion        int64 `json:"chat_messages_region,omitempty"`        /* chat_messages_region integer */
	ChatMessagesSolarsystem   int64 `json:"chat_messages_solarsystem,omitempty"`   /* chat_messages_solarsystem integer */
	ChatMessagesWarfaction    int64 `json:"chat_messages_warfaction,omitempty"`    /* chat_messages_warfaction integer */
	ChatTotalMessageLength    int64 `json:"chat_total_message_length,omitempty"`   /* chat_total_message_length integer */
	DirectTrades              int64 `json:"direct_trades,omitempty"`               /* direct_trades integer */
	FleetBroadcasts           int64 `json:"fleet_broadcasts,omitempty"`            /* fleet_broadcasts integer */
	FleetJoins                int64 `json:"fleet_joins,omitempty"`                 /* fleet_joins integer */
	MailsReceived             int64 `json:"mails_received,omitempty"`              /* mails_received integer */
	MailsSent                 int64 `json:"mails_sent,omitempty"`                  /* mails_sent integer */
}
