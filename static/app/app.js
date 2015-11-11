var nomadApp = angular.module('nomadApp', ['ngResource', 'ngRoute']);

nomadApp.service('client', function(DCs) {
  var exports = {};

  DCs.query(function(data) {
    exports.data = data;
    exports.selectedDC = data[0];
  });

  return exports;
});

nomadApp.controller('mainController', function($scope, $log, $interval,
  UsagePerDC, client) {

  getUsage = function() {
    UsagePerDC.get({
      id: client.selectedDC
    }, function(data) {
      $scope.data = data;
    });
  };

  getUsage();
  $interval(function() {
    getUsage();
  }, 10000);

  $log.log("mainController loaded");
});

nomadApp.controller('nodesController', function($scope, $log, $interval,
  NodesPerDC, client) {

  getUsage = function() {
    NodesPerDC.get({
      id: client.selectedDC
    }, function(data) {
      $scope.data = data;
    });
  };

  getUsage();
  $interval(function() {
    getUsage();
  }, 10000);

  $log.log("nodesController loaded");
});

nomadApp.controller('menuController', function($scope, $location, $log, client,
  Region) {
  $scope.isActive = function(viewLocation) {
    return viewLocation === $location.path();
  };
  $scope.client = client;
  Region.get(function(data) {
    $scope.region = data.region;
  });
  $scope.selectDC = function(dc) {
    $scope.client.selectedDC = dc;
  };

  $log.log("menuController loaded");
});

nomadApp.factory("Region", function($resource) {
  return $resource("/api/region");
});

nomadApp.factory("DCs", function($resource) {
  return $resource("/api/datacenters");
});

nomadApp.factory("UsagePerDC", function($resource) {
  return $resource("/api/usage/:id");
});

nomadApp.factory("NodesPerDC", function($resource) {
  return $resource("/api/usage/:id/nodes");
});

nomadApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
    when('/resources', {
      templateUrl: 'partials/resources.html',
      controller: 'mainController'
    }).
    when('/nodes', {
      templateUrl: 'partials/nodes.html',
      controller: 'nodesController'
    }).
    otherwise({
      redirectTo: '/resources'
    });
  }
]);
