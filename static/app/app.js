var nomadApp = angular.module('nomadApp', ['ngResource', 'ngRoute']);

nomadApp.controller('mainController', function($scope, $log, $interval,
  UsagePerDC) {

  getUsage = function() {
    UsagePerDC.get({
      id: "ams3"
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
  NodesPerDC) {

  getUsage = function() {
    NodesPerDC.get({
      id: "ams3"
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
