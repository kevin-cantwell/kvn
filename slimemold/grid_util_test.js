/**
 * Tests for grid convenience routines / logic.
 *
 * @author Sam Pottinger
 * @license GNU GPL v3
**/

var constants = require("./constants");
var grid_util = require("./grid_util");
var models = require("./models");
var test_util = require("./test_util");


/**
 * Test detection of positions outside the bounds of a grid.
**/
exports.testIsInRange = function(test)
{
    var testGrid = new models.Grid(5, 5);

    var withinPos = new models.GridPosition(1, 1);
    var tooHighPos = new models.GridPosition(1, 6);
    var tooLowPos = new models.GridPosition(1, -1);
    var tooLeftPos = new models.GridPosition(-1, 1);
    var tooRightPos = new models.GridPosition(6, 1);

    test.equal(true, grid_util.isInRange(testGrid, withinPos));
    test.equal(false, grid_util.isInRange(testGrid, tooHighPos));
    test.equal(false, grid_util.isInRange(testGrid, tooLowPos));
    test.equal(false, grid_util.isInRange(testGrid, tooLeftPos));
    test.equal(false, grid_util.isInRange(testGrid, tooRightPos));

    test.done();
}


/**
 * Test getting the positions that neighbor a position.
**/
exports.testGetNeighborPos = function(test)
{
    var testGrid = new models.Grid(5, 5);

    var centerPos = new models.GridPosition(1, 1);
    var upperLeft = new models.GridPosition(0, 0);
    var upperCenter = new models.GridPosition(1, 0);
    var upperRight = new models.GridPosition(2, 0);
    var midLeft = new models.GridPosition(0, 1);
    var midRight = new models.GridPosition(2, 1);
    var bottomLeft = new models.GridPosition(0, 2);
    var bottomCenter = new models.GridPosition(1, 2);
    var bottomRight = new models.GridPosition(2, 2);
    var neighborPositions = grid_util.getNeighborPos(testGrid, centerPos);

    test.ok(!test_util.testObjInArray(centerPos, neighborPositions));
    test.ok(test_util.testObjInArray(upperLeft, neighborPositions));
    test.ok(test_util.testObjInArray(upperCenter, neighborPositions));
    test.ok(test_util.testObjInArray(upperRight, neighborPositions));
    test.ok(test_util.testObjInArray(midLeft, neighborPositions));
    test.ok(test_util.testObjInArray(midRight, neighborPositions));
    test.ok(test_util.testObjInArray(bottomLeft, neighborPositions));
    test.ok(test_util.testObjInArray(bottomCenter, neighborPositions));
    test.ok(test_util.testObjInArray(bottomRight, neighborPositions));

    test.done();
};


/**
 * Test getting cells that neighbor a cell.
**/
exports.testGetNeighbors = function(test)
{
    var testGrid = new models.Grid(5, 5);

    var centerPos = new models.GridPosition(1, 1);
    var centerCell = new models.GridCell(centerPos, 0, 0);
    var upperLeftPos = new models.GridPosition(0, 0);
    var upperLeftCell = new models.GridCell(upperLeftPos, 1, 0);
    var midRightPos = new models.GridPosition(1, 2);
    var midRightCell = new models.GridCell(midRightPos, 2, 1);
    var bottomCenterPos = new models.GridPosition(2, 1);
    var bottomCenterCell = new models.GridCell(bottomCenterPos, 3, 2);

    testGrid.setCellNoChem(centerCell);
    testGrid.setCellNoChem(upperLeftCell);
    testGrid.setCellNoChem(midRightCell);
    testGrid.setCellNoChem(bottomCenterCell);

    var neighbors = grid_util.getNeighbors(testGrid, centerCell);

    test.ok(test_util.testObjInArray(upperLeftCell, neighbors));
    test.ok(test_util.testObjInArray(midRightCell, neighbors));
    test.ok(test_util.testObjInArray(bottomCenterCell, neighbors));
    test.done();
};


/**
 * Test getting the state of a cell given its coordinates.
**/
exports.testGetCellStateCoordiantes = function(test)
{
    var testGrid = new models.Grid(5, 5);
    var pos1 = new models.GridPosition(1, 1);
    var pos2 = new models.GridPosition(2, 2);
    var cell1 = new models.GridCell(pos1, 1, 2);

    testGrid.setCellNoChem(cell1);

    var state = grid_util.getCellStateCoordiantes(
        testGrid, pos1.getX(), pos1.getY());
    test.equal(1, state);

    state = grid_util.getCellStateCoordiantes(
        testGrid, pos2.getX(), pos2.getY());
    test.equal(0, state);

    test.done();
};


/**
 * Test determining if a cell is occupied by an organism given coordinates.
**/
exports.testIsCoordianteOccupiedByOrganism = function(test)
{
    var testGrid = new models.Grid(5, 5);
    var pos1 = new models.GridPosition(1, 1);
    var pos2 = new models.GridPosition(2, 2);
    var cell1 = new models.GridCell(pos1, constants.OCCUPIED_OBSTACLE, 2);
    var cell2 = new models.GridCell(pos2, constants.OCCUPIED_FOOD, 2);

    testGrid.setCellNoChem(cell1);
    testGrid.setCellNoChem(cell2);

    var isOccupied = grid_util.isCoordianteOccupiedByOrganism(
        testGrid, pos1.getX(), pos1.getY());
    test.ok(!isOccupied);

    isOccupied = grid_util.getCellStateCoordiantes(
        testGrid, pos2.getX(), pos2.getY());
    test.ok(isOccupied);

    test.done();
};


/**
 * Test enumerating "active" cells in a grid.
**/
exports.testGetActiveCells = function(test)
{
    var testGrid = new models.Grid(5, 5);
    var pos1 = new models.GridPosition(1, 1);
    var pos2 = new models.GridPosition(2, 1);
    var cell1 = new models.GridCell(pos1, constants.OCCUPIED_ORGANISM, 1);
    var cell2 = new models.GridCell(pos2, constants.OCCUPIED_ORGANISM, 1);

    testGrid.replaceCell(cell1, cell1);
    testGrid.replaceCell(cell2, cell2);

    var activeCells = grid_util.getActiveCells(testGrid);
    var activePositions = activeCells.map(function (x) {
        return x.getPos();
    });

    for(var y=0; y<=2; y++)
    {
        for(var x=0; x<=3; x++)
        {
            var targetPos = new models.GridPosition(x, y);
            test.ok(
                test_util.testObjInArray(targetPos, activePositions)
            );
        }
    }

    test.done();
};
