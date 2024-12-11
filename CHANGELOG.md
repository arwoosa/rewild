# CHANGELOG 1.1.39
## Changes
- Fixes to polaroid (environment variable, participant check, participant polaroid count)
- Show user event scheduled count in invitations

# CHANGELOG 1.1.38
## Changes
- Fixes to collaboration log and achievement

# CHANGELOG 1.1.37
## Changes
- Fix star rule

# CHANGELOG 1.1.36
## Changes
- Achievement changes (filter, star type and addition of eligibility checking)
- Collaborative log route change (load event with rewilding ID where there are achievements)

# CHANGELOG 1.1.35
## Changes
- Check radius, save radius as meter (note: haversine in system returns km)
- Event participant feeling validation
- Collaborative log random count
- Event album link update
- Collaboratinve log polaroid check lat and lng. Use param ?is_check=true
- my/event?has_polaroid=[true,false] Allows retrieval on event with or without polaroids

# CHANGELOG 1.1.34
## Changes
- Change exif lat long library 
- Handle images without lat long
- Removal of checks for image GPS info for achievements

# CHANGELOG 1.1.33
## Changes
- Collaborative log message length limit (20)
- Event join message length limit (50)
- Event accounting max record (100)
- Event announcement per category limit (20)
- Event message board maximum limit (9999)
- Event message board disable create after event ends
- My event pagination changes. When ?page= is not passed, load all

# CHANGELOG 1.1.32
## Changes
- Automatically create user friend suggested on event

# CHANGELOG 1.1.31
## Changes
- /user/{userId}/events?past=false

# CHANGELOG 1.1.30
## Changes
- Rename GET user/:id/event to GET user/:id/events
- GET achievement/places - Bugfix

# CHANGELOG 1.1.29
## Changes
- my/event?past=true allows pagination use ?page=[x]

# CHANGELOG 1.1.28
## Changes
- PUT /event/:eventId allow update to not pass in events_photo or events_photo_cover

# CHANGELOG 1.1.27
## Changes
- POST /event/:eventId reverted to PUT /event/:eventId

# CHANGELOG 1.1.26
## Changes
- GET /event/:eventId/participants show label descriptor

# CHANGELOG 1.1.25
## Changes
- GET /event/:eventId/schedule, POST /event/:eventId/schedule: Fix to day counter 
- PUT /event/:eventId changed to POST /event/:eventId to support image upload
- POST /event/:eventId/accounting insert many at once (delete old record and reinsert)

# CHANGELOG 1.1.24
## Changes
- .env.example LENGTH_REWILDING_NAME, LENGTH_REWILDING_IMAGE, LENGTH_REWILDING_REFERENCE_LINK

# CHANGELOG 1.1.24
## Changes
- Change to achievement object
- GET achievement/places - Retrieve 百岳, 小百岳, 保護區

# CHANGELOG 1.1.23
## Changes
- GET /achievement Load from rewilding event data
- Event polaroids to count image lat lng radius

# CHANGELOG 1.1.22
## Changes
- Event cover image - Rename to utf8

# CHANGELOG 1.1.21
## Changes
- CRUD /event/{eventId}/announcement Use one body structure
- POST rewilding name: 30 chars, images: max 3, reference link: max 3
- GET rewilding-searchText, GET rewilding-searchNearby handle if no places found
- All rewilding create will be run across the achievement database to get 百岳, 小百岳, 保護區 data. Creation of event will automatically add this flag to the data

# CHANGELOG 1.1.20
## Changes
- Event accounting paid by - Validation against participant list
- Create event cover photo 
- Change reference to include link to the cover image
- GET /pocket-list/{pocketListId}/items auto generate URL
- DELETE /pocket-list/{pocketListId}/items - bulk pocket list item delete
- PUT /pocket-list/{pocketListId}/items - bulk pocket list item move

# CHANGELOG 1.1.19
## Changes
- Retrieve achievement based on rewilding achievement type

# CHANGELOG 1.1.18
## Changes
- Add in mega.nz
- Event filter by rewilding
- Event accounting paid by
- Event invitation message length
- String length for UTF8 characters (message board message, event invitation)
- Fix schedule bug
- Retrieve past events event?event_past=1
- Separate announcement from message board endpoints

# CHANGELOG 1.1.17
## Changes
- 1.2a, 1.3 - RewildingSearch: Rename helpers to rewilding-searchText, rewilding-searchNearby, rewilding-autocomplete
- _E1.6 - Pin message board

# CHANGELOG 1.1.16
## Changes
- 1.0 - Add events_rewilding_detail key
- _EVENT - Create and update event missing user ID

# CHANGELOG 1.1.15
## Changes
- Tab1 - Pocket List Name Changes
- _E1.6 - Event Message Board Create/Update Character limit
- _E1 - Event Accounting Create/Update Character validation
- 1.2a, 1.3 - RewildingSearch: Add helpers for rewilding:searchText, rewilding:searchNearby, rewilding:autocomplete, rewilding/place/{placeId}
- 1.1 - Removal of Auth
- 1.4.1 - Addition of title retrieval on rewilding_reference_information
- GENERAL - GetMeta function to handle missing scheme (http/https), fix on exit app when error

# CHANGELOG 1.1.14
## Changes
- 3.1 - Removal of Auth Guard
- 2.1c - Update regex check on event name (alphanumeric, chinese, space)

# CHANGELOG 1.1.13
## Changes
- 4.1 - News

# CHANGELOG 1.1.12
## Changes
- _O1.2 - Get user base camp
- _O1.1, _O1.1a, _O1.1b, _O1.2 - Other user page
- 3.3.1.1.1.1d - Fix to experience

# CHANGELOG 1.1.11
## Changes
- _E1.8, _E1.8a, _E1.8.1 - Event join message
- 1.0a, 1.0b, 1.0c, 1.0d, 1.0e, 1.0f, 1.0g, 1.0h - Event filter date
- _A1a - Pocket list max length
- E1 ~ E1.2.1.1.1.1 - Invitation template and send to user
- 3.4 ~ 3.4c - Achievement routes
- 1.5.1a.2 ~ 1.5.1a.3 - Pocket List add rewilding spot
- 2.0 - GET /my/event : Deleted events still appear

# CHANGELOG 1.1.10
## Changes
- 1.4.1 - POST /rewilding pocket list valid check, check if data is available from google location (dont show data if not available)
- _E1.6 - Pinned message

# CHANGELOG 1.1.9
## Changes
- Experience 3.3.1.1.1.1d ~ 3.3.1.1.1.1e
- Maximum polaroids (3x) 3.3.1 ~ 3.3.1g
- Past event date 3.3.2 ~ 3.3.2a
- General fixes to creating event to use one struct

# CHANGELOG 1.1.8
## Changes
- POST event/{id}/schedule one day event bug fix

# CHANGELOG 1.1.7
## Changes
- POST/PUT event missing event type
- Set events_participant_limit & events_payment_fee to 0
- Fix PUT event/{id} response
- POST event/{id}/schedule request and response changes
- GET event/{id}/schedule response changes
- DELETE event/{id}/schedule implement
- DELETE rewilding/{rewildingId} implement


# CHANGELOG 1.1.6
## Changes
- Changes to notifications

# CHANGELOG 1.1.5
## Changes
- Default type ["national_park", "hiking_area", "campground", "camping_cabin", "park", "playground"] (GET rewilding-search)
- Addition of event_meeting_point_name (POST event | PUT event/{id})
- ⁠Pocket List, Reference link as array (POST rewilding)
- ⁠Adjust event_participants object (GET event | my/event)
- Removal of rewilding_type (GET rewilding)
- Clear unnecessary fmt.Println

# CHANGELOG 1.1.4
## Changes
- Handling of breathing points on participation to event
- Changes to schedule GET and POST structure to be the same
- Changes to Rewilding Create to use same service
- Show participants on user events
- Bugfix pocketlist create rewilding
- Removal of /auth routes as it does not belong to this project

# CHANGELOG 1.1.3
## Changes
- Changes in validation errorlist
- Changes in unauthorized HTTP Code 401 Unauthorised
- Pocket list - Count on move/delete item
- Pocket list - Join on rewilding to get detail
- Pocket List - Location info on rewilding
- Rewilding Search - Water type

# CHANGELOG 1.1.2
## Changes
- Removal of duplicated reference endpoint in rewilding and event in favour for one general references endpoint
- Changes in datatype for Lat and Lng fields
- Removal of POST rewilding-register

# CHANGELOG 1.1.1
## Changes
- Changes in validation messagge
- Changes in all UserAgg to omitempty
- General changes to fit API Spec

# CHANGELOG 1.1.0
## Changes
Event 
- Implement OpenWeather CityID on event registration

Badges
- Helper functions to create badges

Notification
- Helper function to create notifications

Users Created Detail
- Collaborative log polaroid
- Collaborative log album link
- User Event

Change OPTIONS to GET
- rewilding/references
