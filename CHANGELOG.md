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
