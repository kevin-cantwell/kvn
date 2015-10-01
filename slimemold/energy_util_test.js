/**
 * Unit tests for energy state and change calculations.
 *
 * @author Sam Pottinger
 * @license GNU GPL v3
**/

var constants = require("./constants");
var energy_util = require("./energy_util");
var grid_util = require("./grid_util");
var models = require("./models");


/**
 * Routine that tests the calculation of inter-cell energy.
 **/
exports.calculateInterCellEnergy = function(test)
{
    var testGrid = new models.Grid(5, 5);
    
    var leftPos = new models.GridPosition(1, 1);
    var testObstacle = new models.GridCell(leftPos, OCCUPIED_OBSTACLE, 0);
    var rightPos = new models.GridPosition(3, 1);
    var testFood = new models.GridCell(rightPos, OCCUPIED_FOOD, 0);
    var centerPos = new models.GridPosition(2, 1);
    var centerCell = new models.GridCell(centerPos, UNOCCUPIED, 0);

    testGrid.setCell(testObstacle);
    testGrid.setCell(testFood);

    var energy = energy_util.calculateInterCellEnergy(
        testGrid,
        centerCell
    );

    test.equal(energy, -1);

    test.done();
}


/**
 * Routine that tests the calculation of intra-cell energy.
**/
exports.calculateInterCellEnergy = function(test)
{
    var testGrid = new models.Grid(5, 5);
    var energy = 0;
    var leftPos = new models.GridPosition(1, 1);
    var testObstacle = new models.GridCell(
        leftPos,
        constants.OCCUPIED_OBSTACLE,
        0
    );
    var rightPos = new models.GridPosition(3, 1);
    var testFood = new models.GridCell(
        rightPos,
        constants.OCCUPIED_FOOD,
        0
    );
    var centerPos = new models.GridPosition(2, 1);
    var centerCell = new models.GridCell(
        centerPos,
        constants.OCCUPIED_ORGANISM,
        0
    );

    testGrid.setCell(testObstacle);
    testGrid.setCell(testFood);

    energy = energy_util.calculateEnergyIf(
        testGrid,
        centerCell.getPos(),
        centerCell
    );
    test.equal(energy, -2.998);

    testGrid.setCell(centerCell);

    centerCell = new models.GridCell(centerPos, constants.UNOCCUPIED, 0);
    energy = energy_util.calculateEnergyIf(
        testGrid,
        centerCell.getPos(),
        centerCell
    );
    test.equal(energy, 3.182);

    test.done();
}
