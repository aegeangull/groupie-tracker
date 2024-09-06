function initMap(locations) {
    // Make a map
    var latLngBounds = new google.maps.LatLngBounds();
    for (var i = 0; i < locations.length; i++) {
      latLngBounds.extend(locations[i].latLng);
    }
    var map = new google.maps.Map(document.getElementById('map'), {
      zoom: 8,
      center: latLngBounds.getCenter(),
      // Make it cobalt style
      styles: [
        {
          "featureType": "all",
          "elementType": "all",
          "stylers": [
              {
                  "invert_lightness": true
              },
              {
                  "saturation": 10
              },
              {
                  "lightness": 30
              },
              {
                  "gamma": 0.5
              },
              {
                  "hue": "#00aaff"
              }
          ]
        },
        {
            "featureType": "administrative.province",
            "elementType": "geometry.stroke",
            "stylers": [
                {
                    "saturation": "100"
                },
                {
                    "lightness": "27"
                }
            ]
        },
        {
            "featureType": "landscape",
            "elementType": "geometry.fill",
            "stylers": [
                {
                    "color": "#32373c"
                }
            ]
        },
        {
            "featureType": "road.highway",
            "elementType": "geometry.fill",
            "stylers": [
                {
                    "saturation": "100"
                },
                {
                    "lightness": "69"
                },
                {
                    "gamma": "1.40"
                }
            ]
        },
        {
            "featureType": "road.highway",
            "elementType": "labels.text.fill",
            "stylers": [
                {
                    "lightness": "100"
                },
                {
                    "saturation": "100"
                }
            ]
        },
        {
            "featureType": "road.highway.controlled_access",
            "elementType": "labels.icon",
            "stylers": [
                {
                    "saturation": "100"
                }
            ]
        },
        {
            "featureType": "road.arterial",
            "elementType": "geometry.fill",
            "stylers": [
                {
                    "saturation": "43"
                },
                {
                    "lightness": "51"
                }
            ]
        },
        {
            "featureType": "road.arterial",
            "elementType": "labels.text.fill",
            "stylers": [
                {
                    "saturation": "45"
                },
                {
                    "lightness": "19"
                }
            ]
        }
      ]
    });
  
    for (var i = 0; i < locations.length; i++) {  // For each location add marker to map
      var marker = new google.maps.Marker({
        position: locations[i].latLng,
        map: map,
        title: locations[i].name
      });
    }
    map.fitBounds(latLngBounds);
  }

  function getLocations(locationNames, callback) {
    var locations = [];
    var geocoder = new google.maps.Geocoder();
    var remaining = locationNames.length;
    for (var i = 0; i < locationNames.length; i++) {  // For each adress get coordinates
      geocoder.geocode({address: locationNames[i]}, function(results, status) {
        if (status == 'OK') {
          locations.push({
            name: results[0].formatted_address,
            latLng: results[0].geometry.location
          });
        } else {
            console.error('Geocode was not successful: ' + status);
        }
        remaining--;
        if (remaining == 0) {
          callback(locations);
        }
      });
    }
}
