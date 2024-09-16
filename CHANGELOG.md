# CHANGELOG 1.1.12
## Changes
- _O1.2 - Get user base camp
- _O1.1, _O1.1a, _O1.1b, _O1.2 - Other user page

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
