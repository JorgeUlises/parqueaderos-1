'use strict';

angular.module('myapp')
  .config(['$routeProvider', function ($routeProvider) {
    $routeProvider
      .when('/Tipo_propietarios', {
        templateUrl: 'views/Tipo_Propietario/Tipo_propietarios.html',
        controller: 'Tipo_PropietarioController',
        resolve:{
          resolvedTipo_Propietario: ['Tipo_Propietario', function (Tipo_Propietario) {
            return Tipo_Propietario.query();
          }]
        }
      })
    }]);
