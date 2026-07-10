/*******************************************************************************
TITLE:
Stable and waterproof OpenSCAD case by pbtec

DESCRIPTION:

highly scalable case for 3D printing. Try it out!

Optimized for Openscad Customizer. Activate it under view/customizer and play around ;-)

No Support needed to print

- for waterproof cases you can use silicone sealing cord with diameters from 1 to 3mm
- to use also without sealing cord. The groove and ridge gives the housing a high stability and tightness.
- Use of Hot melt copper nuts, regular nuts or square nuts
- define outer vertical radius of corners
- echo output in console shows inner and outer size and more
- echo output shows the needed length of the screws
- use screws from m2 up to m5
- default are 4 screws at each corner. For large cases add addtional ones in the middle of both x and y sides if needed
- use several predefined wall mount holder, some with multiple mounting holes (up to 3)
- the screw holes can now be configured to go through the case or not. The regular nut can also be moved to the bottom of the case
- custumizable 1 to 3 holes on each side for cable glands or similar
- holder for up to 3 pcbs or devices
- bottom mounting holes for wall mounting (not recommended for highly water-resistant housings)

Important!
- If you use standard nuts, you need to pause the printer a certain level to insert the nuts
- If you need a stable and waterproof case please print with 100% infill
- I am aware that sometimes, when using too big or too small parameters, there are some rendering issues.
  To prevent such issues change only one parameter at once and check the result.

for waterproofness see https://blog.prusaprinters.org/watertight-3d-printing-part-2_53638/

This work is licensed under a Creative Commons (4.0 International License)
Attribution-Noncommercial-Share Alike
not allowed | Sharing without ATTRIBUTION
allowed     | Remix Culture allowed
not allowed | Commercial Use
not allowed | Free Cultural Works
not allowed | Meets Open Definition

AUTHOR:
pbtec / pb-tec.ch

use https://paypal.me/pbtec if you want to spend me a coold beer. Thanks in advance :-)

VERSION:
V   KZZ DATE     COMMENT
6.0 pb  31.07.21 First Version to share
6.1 pb  15.12.22 Fixed some problems with german characters (Umlaute) in the code
7.0 pb  15.03.25 Added the option to use hot melt copper nut inserts
7.1 pb  16.03.25 Added solution for 3 pcb holders and 3 side wall holes on each side
7.2 pb  16.03.25 Added separate gasket option to print it with flexible material like TPU
7.3 pb  18.04.25 Fixed again some special character in the comments preventing the code to run
7.4 pb	02.05.25 Crazy, I found unexpected characters (-) on line 36, which turned out to be the issue.
7.5 PhuNguyenPT 08.07.26 Added NutStyle 5: printed internal M3 threads for the 4 corner screw bosses using threads.scad (rcolyer fork), no heat-set insert/nut needed - screws (flat/countersunk head) thread straight into the printed boss. Set as new default NutStyle.
7.6 PhuNguyenPT 08.07.26 Battery holder mounting: replaced the 6 straight-through clearance holes with printed M3 thread bosses (threads.scad), same technique as NutStyle 5. Blind holes only - floor stays sealed on the exterior underside. No nut needed; the battery holder's own thin ~2mm tab is used only as a clearance pass-through.
7.7 PhuNguyenPT 08.07.26 Added PCB1() reference block (45.0 x 40.0 x 10.1mm, 2.5mm corner holes) resting on Device Holder 1's 4 M2.5 standoffs. Display-only visual, same pattern as Breadboard()/BatteryHolder(); hole spacing (32.5 x 37.5mm) already matched the existing DeviceHolder_X_Distance1/DeviceHolder_y_Distance1 values.
7.8 PhuNguyenPT 08.07.26 Device Holder 1 (solar PCB controller) mounting bosses: added optional printed internal M2.5 threads (threads.scad ScrewThread(), same technique as NutStyle 5 / battery holder bosses) to the DeviceHolder() module, enabled via new DeviceHolder1_UseThread/DeviceHolder1_ScrewSize/DeviceHolder1_ThreadPitch settings. Screw threads straight into the boss, no nut/insert needed. DeviceHolder2/3 unaffected (default to the original plain clearance hole).
7.9 PhuNguyenPT 08.07.26 Fixed Device Holder 1 (solar PCB controller) M2.5 printed thread bosses:
    the ScrewThread() cut previously fell 0.01mm short of BOTH the top and bottom faces
    (height was CylHeight-0.02, symmetric margin), fully sealing the thread cavity inside
    solid plastic with no screw entry. Changed height to CylHeight so the cut now overshoots
    the top (entry) face by 0.01mm, matching the already-working NutStyle 5 corner boss and
    battery holder boss convention. Base still pulls back 0.01mm (blind, sits on the floor).
7.10 PhuNguyenPT 10.07.26 Fixed "2 non-manifold edges" STL export error that only appeared when
    EnableMountHolder=true. Root cause: roundedBox() built its rounded box out of a union of
    overlapping cubes + edge cylinders + corner spheres, all tangent to each other at exact
    mathematical boundaries - a classic source of coincident/degenerate faces once unioned into
    the wall mount holder (all 5 MountHolderStyle variants use it) and then differenced with
    cutting cubes/holes.
7.11 PhuNguyenPT 10.07.26 Switched roundedBox() to BOSL2's cuboid(rounding=...) (see
    include <BOSL2/std.scad> near the top of the file) instead of a hand-rolled hull() - same
    idea (a single coherent construction instead of a union of tangent primitives), but using
    the well-tested library implementation as requested. Two things had to be fixed for BOSL2 to
    drop in cleanly: (1) the BOSL2 line must be "include", not "use" - "use" only imports
    modules/functions in OpenSCAD, not variables, so CENTER/UP/DOWN/etc. would come back
    undefined. (2) BOSL2's cylinder() overrides the built-in one with stricter validation
    (h must be > 0); the pre-existing ScrewCut(m,h,v) corner-screw-cut calls all pass v=0 ("no
    extra countersink depth"), which built a zero-height cylinder(h=v,...) that the old built-in
    cylinder() silently no-op'd but BOSL2's version correctly rejects. Fixed by guarding that
    cylinder with "if (v>0)" in ScrewCut() - a v=0 sink depth is nothing to cut anyway, so this
    is the correct behavior regardless of which cylinder() implementation is active. Call
    signature and resulting shape of roundedBox() are unchanged - no edits needed at any
    MountHolder() call site.
7.12 PhuNguyenPT 10.07.26 Circularity audit: searched the whole file for every remaining
    "$fn=30" override (the battery-holder mounting boss cylinders were already fixed above) and
    removed all of them so those cylinders inherit the global $fn=80 instead:
    - BatteryHolder() reference block: the 6 tab bore holes + countersink cones (display-only,
      not part of the printed geometry, but was visibly more faceted than the rest of the model).
    - PCB1() reference block: the 4 corner mounting holes (same, display-only).
    Left untouched, on purpose: $fn=40/60 inside MountHolder() (styles 1/2/4/5) - local
    resolution choices scoped to that module, not a leftover/inconsistency; and the hex-nut
    cylinder's $fn=6 (NutCutSquare/hex helper around line 1335) - 6 sides is the actual intended
    hexagon shape for a nut, not a resolution shortcut, so it must stay exactly 6.
*******************************************************************************/

// Requires threads.scad (Ryan A. Colyer's library, CC0):
// https://github.com/rcolyer/threads-scad
// Download threads.scad and place it in the same folder as this file (or in
// your OpenSCAD library folder). We use its ScrewThread() primitive directly
// (the same subtraction ScrewHole() performs internally) to cut real printed
// internal M3 threads into the 4 corner screw bosses (NutStyle 5 below),
// instead of a plain hole for a heat-set insert/self-tapper.
use <threads-scad/threads.scad>

// Requires BOSL2 (Belfry OpenSCAD Library v2, https://github.com/BelfryCAD/BOSL2, BSL-2 license):
// Download/clone the repo and place the whole folder (rename to "BOSL2") next to this file, or in
// your OpenSCAD library folder. NOTE: this is "include", not "use" - BOSL2 defines constants like
// CENTER/UP/DOWN/LEFT/RIGHT as plain variables, and OpenSCAD's "use" statement only imports
// modules/functions, not variables, so "use" would leave CENTER undefined. BOSL2's cuboid()
// replaces the old hand-rolled roundedBox() (see module below) - it builds rounded boxes as a
// topologically clean solid internally, which is why it does not produce the non-manifold edges
// that the old cube+cylinder+sphere union version produced once the wall mount holder (5 styles,
// EnableMountHolder) unions/differences several rounded boxes together.
include <BOSL2/std.scad>

/* [Render quality settings] */
// Set to at least to 150 before render and save as .stl file, otherwise you can go down to 40 for quick 3D view
$fn                       = 80;   // [20:1:300]
NozzleDia = 0.2;
ClearanceGap = 0.75*NozzleDia;
/* [View settings] */
// Shows the Bottom of the case
ShowBottom                = true;
// Shows the top of the case
ShowTop                   = true;
// Distance between top and bottom (if both are side by side displayed)
DistanceBetweenObjects    = 10;
// Shows the housing assembled
ShowCaseAssembled         = false;
// Shows the gasket
ShowGasket                = true;

/* [Control cuts (use only one at a time)] */
// To see the nuts inside (best view if not assembled showed)
SeeNutCut                 = false;
// To see the groove, ridge and Screw (best view if assembled showed)
SeeGrooveRidgeScrew       = false;

/* [Case settings] */
// Length of the case
Caselength                = 180;
// Width of the case - widened to fit the breadboard plus a battery holder side by side
CaseWidth                 = 187;
// Height of the case
CaseHeight                = 60;
// Splitt the Case height into bottom and top, check for the needed screws in echo output (console)
CutFromTop                = 30.0;
// Thickness for the bottom and top wall (vertical walls needs to be calculated)
BottomTopThickness        = 3.0;
// If this is bigger than the needed cylinder around the screw it will be ignored
CaseRadius                = 12.0;

/* [Breadboard reference block] */
// Show a reference rectangle representing the breadboard (visual only, it is not merged into the printed bottom/top parts)
ShowBreadboard            = true;
// Length of the breadboard (X direction)
BreadboardLength          = 165;
// Width of the breadboard (Y direction)
BreadboardWidth           = 105;
// Height of the breadboard (Z direction / vertical clearance needed)
BreadboardHeight          = 10;
// Move the breadboard left/right from the case center
BreadboardOffset_X        = 0;
// Move the breadboard forward/back from the case center - shifted toward one side to leave room for the battery holder
BreadboardOffset_Y        = 32.0;
// Lift the breadboard up from the inner floor (0 = resting directly on the floor)
BreadboardOffset_Z        = 0;

/* [Battery holder reference block] */
// Show a reference rectangle representing the battery holder (visual only, it is not merged into the printed bottom/top parts)
ShowBatteryHolder         = true;
// Length of the battery holder (X direction)
BatteryHolderLength       = 79;
// Width of the battery holder (Y direction)
BatteryHolderWidth        = 59.5;
// Height of the battery holder (Z direction)
BatteryHolderHeight       = 22;
// Move the battery holder left/right from the case center - positioned close to a corner
BatteryHolderOffset_X     = -42.5;
// Gap between the breadboard's near edge and the battery holder's near edge, along Y.
// Tune THIS to widen/narrow the clearance between the two - BatteryHolderOffset_Y below is
// derived from it (and BreadboardOffset_Y/BreadboardWidth, both already set above), so the
// battery holder AND its printed mounting bosses in BodyBottom (which reuse this same offset)
// move together automatically. Default reproduces the original hand-placed 0.75mm gap.
BreadboardBatteryHolderGapY = 3.0;   // [0:0.25:50]
// Move the battery holder forward/back from the case center - derived from the gap above so
// it stays anchored to the breadboard's edge instead of a fixed absolute coordinate.
BatteryHolderOffset_Y     = (BreadboardOffset_Y - BreadboardWidth/2) - BatteryHolderWidth/2 - BreadboardBatteryHolderGapY;
// Lift the battery holder up from the inner floor (0 = resting directly on the floor)
BatteryHolderOffset_Z     = 0;

BatteryHolderHoleDiameter = 3.5;
// Center distances based on your edge measurements (+1.75mm radius)
BatteryHolderHoleGapFromWidthEdge = 11.75;  // 10mm edge gap + 1.75mm radius
BatteryHolderHoleGapFromLengthEdge = 10.25; // 8.5mm edge gap + 1.75mm radius
BatteryHolderHoleSpacing = 19.5;            // 16mm edge spacing + 3.5mm (two radii)

/* [Battery holder mounting bosses (printed M3 threads, threads.scad, no nut needed)] */
// Height of the printed thread boss rising from the inner floor at each of the 6 battery holder mounting points.
// This is a BLIND hole (does not go through the floor to the outside), so the case floor stays sealed.
// The battery holder's own mounting tab (only ~2mm thick) is too thin to hold an M3 thread reliably on its own -
// the screw passes through it as a plain clearance hole and threads into this boss instead.
// Sized for an M3 x 12mm flat/countersunk head screw: the countersunk head sinks flush into the
// battery holder's ~2mm tab, leaving 12-2 = 10mm of shank below the tab that needs to thread in.
// +2mm added on top of that as bottom-of-hole clearance, so the screw tip never bottoms out on solid
// plastic before the countersunk head is fully seated in the tab (print tolerance / thread lead-in margin).
BatteryHolderBossHeight     = 6.0;  // [2:0.1:20]
// Outer diameter of each printed thread boss (needs enough wall around the M3 thread cut for strength)
BatteryHolderBossDiameter  = 8.0;  // [5:0.1:15]

/* [Battery holder reference model: floor tab, walls & slots (illustrative only, display block)] */
// Straight bore depth at each mounting hole - bottom of the tab, exits toward the boss
BatteryHolderTabBoreDepth   = 2.0;  // [0.5:0.1:5]
// Countersunk (flat, rounded) taper depth above the bore (set to 0 for flat button screws)
BatteryHolderTabCsinkDepth  = 0.0;  // [0.0:0.1:5]
// Total tab/floor thickness (bore + countersink)
BatteryHolderFloorThickness = BatteryHolderTabBoreDepth + BatteryHolderTabCsinkDepth;
// Countersink's wider opening radius: +1mm (0.1cm) over the mounting hole's own radius
BatteryHolderCsinkTopRadius = BatteryHolderHoleDiameter/2 + 1.0;
// Outer wall thickness in the WIDTH direction (the long side walls running along the length axis)
BatteryHolderWidthWallThickness  = 2.0;   // [0.5:0.1:5]
// Outer wall thickness in the LENGTH direction (the end walls running along the width axis)
BatteryHolderLengthWallThickness = 1.0;   // [0.5:0.1:5]
// Height of the width-direction side walls (and the 2 internal slot dividers, which run the same direction)
BatteryHolderWidthWallHeight     = 12.0;  // [5:0.1:30]
// Height of the length-direction end walls
BatteryHolderLengthWallHeight    = 21.5;  // [5:0.1:30]
// Thickness of the 2 internal divider walls that split the interior into 3 equal battery slots
BatteryHolderSlotWallThickness   = 1.0;   // [0.5:0.1:5]

/* [Groove settings] */
// If using a SealingCord use the SealingCord diameter, otherwise x-times of your 3D Printer Nozzle (0.8/1.2/...) --> Ridge gets perfect for printing
GrooveWidth               = 1.2;   // [0.8:0.1:3]
// Not less than 1mm (for stability) and not more than 3mm --> Best 2mm
GrooveDepth               = 2.0;   // [1.0:0.1:3]
// Space between Groove and Ridge for a perfect fit, usualy 0.2 or 0.3 for FDM depending on your printer quality
Space                     = 0.3;   // [0.0:0.1:0.4]
// Addtional vertical room for the pressed sealing Cord. For sealing cord 1.5mm -->0.5 // for 2mm -->0.8 // for 2.5 -->1.0 // If no sealing cord then set this parameter to 0.
AddGrooveDepthForSealing  = 0.8;   // [0.0:0.1:3]
GasketSpace       = Space + ClearanceGap;      // looser than the rigid ridge fit, tied to nozzle-derived clearance
GasketCompression = 0.3;  // [0.1:0.05:0.5]
// Range Inside groove/ridge. Usualy 2 times or more the printer nozzle. For best stability at least 0.8
InnerBorder               = 0.8;   // [0.8:0.1:4]
// Range Outside groove/ridge . Usualy 2 times or more the printer nozzle. For best stability at least 0.8
OuterBorder               = 0.8;   // [0.8:0.1:4]

/* [Case Screw settings] */
//2=m2/2.5=m2.5/3=m3/4=m4/5=m5   // max m5, larger sizes do not fit
MetricScrewSize          = 3.0;     // [2:m2, 2.5: m2.5, 3: m3, 4: m4, 5: m5]
// Chose your Screw head
ScrewHeadType            = 1;      // [1:1 Countersunk head screw, 2: 2 Round or hex screw head - Counterbore, 3: 3 Exposed head - No counter]
// for round or hex screw head select the height of the head
ScrewHeadHeight  = 2.6;   // [0.0:0.1:10]
// for round head screw select the diameter of the head, for hex screw select the diameter size of your Socket wrench diameter
ScrewHeadDiameter  = 6;   // [0.0:0.1:20]
// Screw or hot melt nut Hole deepness (deepness in the Body/Cylinder). If too big, through hole possible
HoleDeepness              = 20.0 ; //[1:0.1:500]
// Adds additional Screws on X axis (for large cases) --> Try it out
XAdditionalScrew          = false;  // can be true or false / Adds additional Screws on X axis (for large cases) --> Try it out
// Adds additional Screws on Y axis (for large cases) --> Try it out
YAdditionalScrew          = false;  // can be true or false / Adds additional Screws on Y axis (for large cases) --> Try it out

/* [Nut general settings] */
NutStyle          = 5;      // [1:1 Hot melt copper nut > New from Version 7.0, 2: 2 Standard nuts > Pause during printing to insert standard nuts, 3: 3 Square nuts > Pause during printing to insert square nuts, 4 : 4 Square nuts with external insert > No need for pause during print - but outer holes are visible, 5: 5 Printed internal threads (threads.scad) > No insert/nut needed, screw threads directly into the printed corner boss]
// Size of material (plastic) above nut/square nut (3mm or better more). Use for NutStyle 2,3 and 4). The more, the more stable but need longer screws. Here you can also move the nut to the bottom of the case (since Version 7.0)
NutSink                   = 4.0; // [0:0.1:500]


/* [Hot melt copper nut settings (use for NutStyle 1) can also be used for self-tapping screws] */
// Hole diameter for self-tapped screw or hot melt copper nut (Measure the diameter))
HoleDiaThread             = 3.2 ; //[1:0.1:10]
// Length of the hot melt copper nut (Measure the length) - It is just used to calculate the min length of the screw to output in console
HolelengthHotMeltNut      = 10 ; //[1:0.1:30]

/* [Printed internal thread settings (use for NutStyle 5, threads.scad, corner screws)] */
// Thread pitch in mm. ISO metric coarse pitch: m2=0.4 // m2.5=0.45 // m3=0.5 // m4=0.7 // m5=0.8
// (threads.scad's own ThreadPitch(3) lookup already returns 0.5 for M3, kept explicit here for clarity/customizer control)
ThreadPitch               = 0.5;
// Thread flank angle in degrees, passed to threads.scad as tooth_angle (30 = standard ISO metric 60 deg thread profile - do not change unless you know what you are doing)
ThreadAngle               = 30;
// How deep the printed thread engages into the corner boss (reuses HoleDeepness by default, edit HoleDeepness above to change both together)
ThreadEngagementLength    = HoleDeepness;
// Passed straight through to threads.scad's "tolerance" parameter on ScrewThread(). This is the FDM print-fit compensation:
// it slightly enlarges/loosens the cut thread so a real M3 screw actually threads into the printed hole.
// Derived the same way as ClearanceGap above: 2 * (0.75*NozzleDia). threads.scad's own default is 0.4; increase toward that
// (or higher) if the printed thread comes out too tight, decrease if too loose/stripped.
ThreadFitComp             = 2*ClearanceGap;

/* [Standard nut settings (use for NutStyle 2)] */
// Nut Settings / As there are (or I have) many different nuts dimensions, the size must be specified / Do not add separation space, only the real measurement // m2=1.5 // m2.5=1.9 // m3=2.4 // m4=3.2 // n5=3.8
NutHigh                   = 2.4;
// Distance between the paralell sides / Do not add separation space, only the real measurement // m2=3.8 // m2.5=4.9 // m3=5.4 // m4=6.9 // m5=7.9
NutDia                    = 5.4;

/* [Square Nut settings (Use for NutStyle 3 + 4)] */
// Select the high of the square nut / Do not add separation space, only the real measurement
SquareNutHigh             = 1.9;
// Select the size of the square nut / Do not add separation space, only the real measurement
SquareNutSize             = 5.4;
// Square nut insert from which side (NutStyle 4 only)
EdgeSquareNutInsertFrom_X = true;

/* [Wall mount holder settings] */
// Select if you need a mount holder
EnableMountHolder         = true;
// Chose your desired wall mount style
MountHolderStyle          = 5;      // [1:Style 1, 2: Style 2, 3: Style 3, 4 : Style 4, 5 : Style 5]
// Some styles (1-3) allow more than one hole
CountOfMountHolderHoles   = 1;     // [1:One hole centered, 2: Two holes, 3: Three holes]
MountHolderHoleDiameter   = 5;   // [1:0.1:10]
MountHolderThickness      = 4.0;   // [2:0.1:10]

/* [PCB/Device holder 1 settings] */
// Activate customizable PCB/Device holder 1
ShowDeviceHolder1         = true;

// Hole in the cylinder for the screw (Fits standard M2.5 screws) - only used as a plain clearance
// hole if DeviceHolder1_UseThread below is false
ScrewHoleDiameter1        = 2.2;

// Use printed internal M2.5 threads (threads.scad ScrewThread(), same technique as the corner case
// screws / NutStyle 5 and the battery holder mounting bosses) instead of a plain clearance hole.
// The screw (e.g. flat/countersunk M2.5) threads straight into the printed boss - no nut needed.
DeviceHolder1_UseThread    = true;
// Nominal screw size for the printed thread on Device Holder 1 (M2.5)
DeviceHolder1_ScrewSize    = 2.5;
// ISO metric coarse thread pitch in mm for the size above: m2=0.4 // m2.5=0.45 // m3=0.5 // m4=0.7 // m5=0.8
DeviceHolder1_ThreadPitch  = 0.45;

// The diameter of the screw cylinder (standoff width)
ScrewCylinderDiameter1    = 6.0;

// The height of the screw cylinders (standoff height)
ScrewCylinderHeight1      = 6.0;

// Distance between the holes in X direction (4.0cm side width)
DeviceHolder_X_Distance1  = 32.5;

// Distance between the holes in Y direction (4.5cm side length facing Side A)
DeviceHolder_y_Distance1  = 37.5;
// Clearance between Device Holder 1's outer edge and the inner face of Wall A (+X side).
// Tweak THIS value until the holder sits where you want next to the wall - Offset_X_1 below
// is now derived from it, so the holder (and PCB1, which reuses Offset_X_1) automatically
// stays anchored to Wall A if Caselength changes, instead of needing to be re-tuned by hand.
DeviceHolder1_ClearanceToWallX = 30.95;   // [0:0.05:150]

// Clearance between Device Holder 1's outer edge and the inner face of the -Y wall.
// Same idea as above - Offset_Y_1 below is derived from it, keeping the holder anchored
// to that wall if CaseWidth changes.
DeviceHolder1_ClearanceToWallY = 16.95;   // [0:0.05:150]

// Derived: X offset from case center. Computed from half the case length, minus the side
// wall thickness, minus the clearance-to-wall above, minus half the hole spacing - this
// reproduces the original hand-tuned 40mm at the current Caselength (180mm), but now moves
// with Caselength automatically.
// NOTE: SideWallThickness isn't assigned until later in the file (Calculated settings block),
// and OpenSCAD evaluates top-level assignments in file order rather than as a lazy dependency
// graph, so its definition (InnerBorder+GrooveWidth+OuterBorder) is inlined here instead of
// referencing the name directly - both InnerBorder and GrooveWidth/OuterBorder are already
// defined earlier in the file, above this point.
Offset_X_1                = Caselength/2 - (InnerBorder+GrooveWidth+OuterBorder) - DeviceHolder1_ClearanceToWallX - DeviceHolder_X_Distance1/2;

// Derived: Y offset from case center, same technique, anchored toward the -Y wall.
// Reproduces the original hand-tuned -55mm at the current CaseWidth (187mm), but now moves
// with CaseWidth automatically. Same SideWallThickness-inlining note as above applies.
Offset_Y_1                = -(CaseWidth/2 - (InnerBorder+GrooveWidth+OuterBorder)) + DeviceHolder1_ClearanceToWallY + DeviceHolder_y_Distance1/2;

/* [PCB 1 reference block (visual only, sits on the 4 M2.5 standoffs above)] */
// Show a reference block for the PCB resting on Device Holder 1's 4 corner standoffs
ShowPCB1                  = true;
// Overall PCB length (45.0mm official dimension) - rendered along Y, matching DeviceHolder_y_Distance1
PCB1_Length               = 45.0;
// Overall PCB width (40.0mm official dimension) - rendered along X, matching DeviceHolder_X_Distance1
PCB1_Width                = 40.0;
// Overall PCB height/clearance, Z direction (10.1mm official dimension)
PCB1_Height               = 10.1;
// Diameter of the 4 corner mounting holes (fits M2.5 screws)
PCB1_HoleDiameter         = 2.5;

/* [PCB/Device holder 2 settings] */
// Activate customizable PCB/Device holder 2
ShowDeviceHolder2         = false;
// Hole in the cylinder for the screw // 2.9 Screw = 2mm hole
ScrewHoleDiameter2        = 2.6;
// The diamter of the screw cylinder
ScrewCylinderDiameter2    = 7;
// The height of the screw cylinders (also the deepness for the screw hole)
ScrewCylinderHeight2      = 10;
// Distance between the holders in X direction
DeviceHolder_X_Distance2  = 15;
// Distance between the holders in Y direction
DeviceHolder_y_Distance2  = 50;
// Move all holders 2 in X direction
Offset_X_2                = -15;
// Move all holders 2 in Y direction
Offset_Y_2                = -10;

/* [PCB/Device holder 3 settings] */
// Activate customizable PCB/Device holder 3
ShowDeviceHolder3         = false;
// Hole in the cylinder for the screw // 2.9 Screw = 2mm hole
ScrewHoleDiameter3        = 3;
// The diamter of the screw cylinder
ScrewCylinderDiameter3    = 9;
// The height of the screw cylinders (also the deepness for the screw hole)
ScrewCylinderHeight3      = 8.0;
// Distance between the holders in X direction
DeviceHolder_X_Distance3  = 40;
// Distance between the holders in Y direction
DeviceHolder_y_Distance3  = 10;
// Move all holders 3 in X direction
Offset_X_3                = 0;
// Move all holders 3 in Y direction
Offset_Y_3                = 0;

/* [Base bottom mounting holes settings. Not recommended for highly water-resistant housings] */
// Activate customizable base bottom mounting holes
ShowBottomMountingHoles     = false;
// Diameter for the mounting holes
BottomMountingHolesDiameter = 3;
// Distance between the mounting holes in X direction
BottomMountingHoles_X_Distance = 40;
// Distance between the mounting holes in Y direction
BottomMountingHoles_y_Distance = 70;
// Move all mounting holes in X direction
BottomMountingHolesOffset_X    = 50;
// Move all mounting holes in Y direction
BottomMountingHolesOffset_Y    = 30;

/* [Wall Holes settings side A (for cable gland cut)] */
// Activate customizable holes for cable gland or similar
ShowSideWallHoles_A        = true;
// Count of holes, if there is an additional screw on X or Y side the hole in the middle is not showed
CountOfSideWallHoles_A     = 3;     //[1:1:3]
// Diameter of the holes
SideWallHoleDiameter_A = 12.5 + 2*ClearanceGap;    //[1:0.1:80]
// PG7 nominal 12.5mm + FDM clearance, not hardcoded
// Add or decrease height position (up and down, 0 = centered)
SideWallHoleOffset_Z_A     = 5;
// Add or decrease distance between the holes
SideWallHoleDistance_A     = 20;
// Add or decrease horizontal position of the holes (0 = centered)
SideWallHolePosition_A     = -45;

/* [Wall Holes settings side B (for cable gland cut)] */
// Activate customizable holes for cable gland or similar
ShowSideWallHoles_B        = false;
// Count of holes, if there is an additional screw on X or Y side the hole in the middle is not showed
CountOfSideWallHoles_B     = 1;     //[1:1:3]
// Diameter of the holes
SideWallHoleDiameter_B     = 12.5;  //[1:0.1:80]
// Add or decrease height position (up and down, 0 = centered)
SideWallHoleOffset_Z_B     = 0;
// Add or decrease distance between the holes
SideWallHoleDistance_B     = 20;
// Add or decrease horizontal position of the holes (0 = centered)
SideWallHolePosition_B     = 0;

/* [Wall Holes settings side C (for cable gland cut)] */
// Activate customizable holes for cable gland or similar
ShowSideWallHoles_C        = false;
// Count of holes, if there is an additional screw on X or Y side the hole in the middle is not showed
CountOfSideWallHoles_C     = 1;     //[1:1:3]
// Diameter of the holes
SideWallHoleDiameter_C     = 16.5;  //[1:0.1:80]
// Add or decrease height position (up and down, 0 = centered)
SideWallHoleOffset_Z_C     = 0;
// Add or decrease distance between the holes
SideWallHoleDistance_C     = 20;
// Add or decrease horizontal position of the holes (0 = centered)
SideWallHolePosition_C     = 0;

/* [Wall Holes settings side D (for cable gland cut)] */
// Activate customizable holes for cable gland or similar
ShowSideWallHoles_D        = false;
// Count of holes, if there is an additional screw on X or Y side the hole in the middle is not showed
CountOfSideWallHoles_D     = 1;     //[1:1:3]
// Diameter of the holes
SideWallHoleDiameter_D     = 16.5;  //[1:0.1:80]
// Add or decrease height position (up and down, 0 = centered)
SideWallHoleOffset_Z_D     = 0;
// Add or decrease distance between the holes
SideWallHoleDistance_D     = 20;
// Add or decrease horizontal position of the holes (0 = centered)
SideWallHolePosition_D     = 0;

/* [RJ45 Ethernet port cutouts (positioned in the top piece, above the breadboard)] */
// Activate ethernet port cutout 1
ShowEthernetPort1          = true;
// Activate ethernet port cutout 2
ShowEthernetPort2          = true;
// Width of the rectangular cutout as seen from outside the case
EthernetPortWidth          = 19 + 2*ClearanceGap;
// Height of the rectangular cutout as seen from outside the case
EthernetPortHeight         = 17 + 2*ClearanceGap;
// How far the RJ45 module body extends into the inside of the case, measured from the inner wall face
EthernetPortInsideDepth    = 10 + 2*ClearanceGap;
// Which wall port 1 sits in: A = right (+X), B = front (-Y), C = left (-X), D = back (+Y)
EthernetPort1_Wall         = "C"; // [A,B,C,D]
// Which wall port 2 sits in: A = right (+X), B = front (-Y), C = left (-X), D = back (+Y)
EthernetPort2_Wall         = "C"; // [A,B,C,D]
// Horizontal position of port 1 along its wall (0 = centered)
EthernetPort1_Position     = 13.5;
// Horizontal position of port 2 along its wall (0 = centered)
EthernetPort2_Position     = 50.5;
// Fine-tune vertical position of port 1, relative to the default (centered in the top piece, above the breadboard). 0 = default
EthernetPort1_OffsetZ      = 0;
// Fine-tune vertical position of port 2, relative to the default (centered in the top piece, above the breadboard). 0 = default
EthernetPort2_OffsetZ      = 0;
// Add a friction-fit collar (thickened sleeve) around each port that extends inward for a snug grip on the RJ45 module body
ShowEthernetPortCollar     = true;
// Thickness of added material around the port opening, on each side (needs a real friction-fit channel, not just a clip)
EthernetPortCollarWall     = 1.6;
// Show a reference model of the physical RJ45 jack module in each port (visual only, not part of the printed model)
ShowRJ45Model              = true;
// =========================  C A L C U L A T E D   S E T T I N G S (do not change!!!) ===================================

// Calculated Screw settings (do not change!!!)
ScrewHoleDia              = MetricScrewSize+1;
ScrewHeadDia              = MetricScrewSize*2;
ScrewCountersink          = (MetricScrewSize+8)/14-0.5;

// Calculated settings for Ridge (do not change!!!)
RidgeHeight               = GrooveDepth-Space;
RidgeWidth                = GrooveWidth-Space;

GasketRidgeWidth  = GrooveWidth - GasketSpace;  // 1.2 - 0.45 = 0.75mm (rigid ridge is 0.9mm — sits slightly recessed atop it)
// Fraction of the gasket's printed height that gets squeezed out when the screws pull the case shut
GasketHeadroom = AddGrooveDepthForSealing + Space;      // 0.8 + 0.3 = 1.1mm, vertical room above the ridge tip
GasketHeight   = GasketHeadroom / (1 - GasketCompression); // 1.1 / 0.7 ≈ 1.57mm printed height

// Calculated settings for case (do not change!!!)
SideWallThickness             = InnerBorder+GrooveWidth+OuterBorder;
CaseRoundingRadius        = ScrewHoleDia/2+InnerBorder+GrooveWidth+OuterBorder;

// Calculated default vertical position for the ethernet ports (do not change!!!)
// Centered within the top piece's wall height, i.e. above the split line and comfortably clear of the breadboard
EthernetPortDefaultZ      = (CaseHeight-CutFromTop) + CutFromTop/2;
ScrewCornerPos            = [Caselength/2-CaseRoundingRadius,CaseWidth/2-CaseRoundingRadius,0];
ScrewAddXPos              = [0,CaseWidth/2-CaseRoundingRadius,0];
ScrewAddYPos              = [Caselength/2-CaseRoundingRadius,0,0];

// Calculated settings for wall mount holder
MountHolderLenght         = MountHolderHoleDiameter*3;

// if both objects showed
X_ObjectPosition = ((ShowBottom)&&(ShowTop)&&(!ShowCaseAssembled)) ? Caselength/2+DistanceBetweenObjects/2:0;

// If the case is assembled showed
Y_TopRotation = ShowCaseAssembled ? 180:0;
Z_TopHigh = ShowCaseAssembled ? CaseHeight:0;

ShowSizes(); // Show the calculated sizes

//===============================================================================
//                                    M A I N
//===============================================================================




// --> Show the bottom of the case
translate([X_ObjectPosition,0,0]) rotate([0,0,0]) difference(){
    union(){
        BodyBottom();
        // **** Add your bottom case additions here ****
        //cylinder(h=20,d=15,center = true); // Example
    }
    // **** Add your bottom case cuts here ****
    //cylinder(h=15,d=20,center = true); // Example
}

// --> Show the top of the case
translate([-X_ObjectPosition,0,Z_TopHigh+0.03]) rotate([0,Y_TopRotation,0]) difference(){
    union(){
        BodyTop();
        // **** Add your bottom top additions here ****
        //cylinder(h=18,d=10,center = true); // Example
    }
    // **** Add your top case cuts here ****
    //cylinder(h=20,d=5,center = true); // Example
}

// --> Show RJ45 jack reference models at the ethernet port locations (visual only, not part of
// the printed model). Wrapped in the same translate/rotate as the top piece above so it tracks
// correctly whether the case is shown exploded, side-by-side, or assembled (ShowCaseAssembled).
if (ShowTop) translate([-X_ObjectPosition,0,Z_TopHigh+0.03]) rotate([0,Y_TopRotation,0])
{
    if (ShowEthernetPort1) RJ45Model(EthernetPort1_Wall, EthernetPort1_Position, CaseHeight-(EthernetPortDefaultZ+EthernetPort1_OffsetZ));
    if (ShowEthernetPort2) RJ45Model(EthernetPort2_Wall, EthernetPort2_Position, CaseHeight-(EthernetPortDefaultZ+EthernetPort2_OffsetZ));
}

if (ShowGasket)
{
    translate([-Caselength/2-Caselength-DistanceBetweenObjects,0,0]) Gasket();
}

// --> Show the breadboard reference block (visual only, not part of the printed model)
// Only shown alongside the bottom case, since it rests on bosses that live in BodyBottom().
if (ShowBottom) translate([X_ObjectPosition,0,0]) Breadboard();

// --> Show the battery holder reference block (visual only, not part of the printed model)
// Only shown alongside the bottom case, since it rests on bosses that live in BodyBottom().
if (ShowBottom) translate([X_ObjectPosition,0,0]) BatteryHolder();

// --> Show the PCB 1 reference block, resting on the 4 M2.5 standoffs (visual only, not part of the printed model)
// Only shown alongside the bottom case, since it rests on standoffs that live in BodyBottom().
if (ShowBottom) translate([X_ObjectPosition,0,0]) PCB1();

//===============================================================================
//                                  M O D U L E S
//===============================================================================



module BodyBottom () {
    if(ShowBottom)
    {
        // Battery holder mounting boss/hole reference positions (used by both the boss additions
        // in the union() below and the printed-thread cuts in the difference() below)
        BatteryHolder_LeftX = BatteryHolderOffset_X - (BatteryHolderLength/2);
        BatteryHolder_RightX = BatteryHolderOffset_X + (BatteryHolderLength/2);
        BatteryHolder_BottomY = BatteryHolderOffset_Y - (BatteryHolderWidth/2);

        difference(){
            union()
            {
                rotate([  0,  0,  0]) BodyQuarterBottom(Caselength,CaseWidth,CaseHeight-CutFromTop,CaseRoundingRadius,SideWallThickness);
                rotate([  0,  0,180]) BodyQuarterBottom(Caselength,CaseWidth,CaseHeight-CutFromTop,CaseRoundingRadius,SideWallThickness);
                mirror([  0,  1,  0]) BodyQuarterBottom(Caselength,CaseWidth,CaseHeight-CutFromTop,CaseRoundingRadius,SideWallThickness);
                mirror([  1,  0  ,0]) BodyQuarterBottom(Caselength,CaseWidth,CaseHeight-CutFromTop,CaseRoundingRadius,SideWallThickness);

                // Add here additional parts

                if (EnableMountHolder)
                {
                    color("SteelBlue") if (MountHolderStyle!=5)
                    {
                    translate([0,CaseWidth/2,0]) MountHolder(MountHolderThickness,MountHolderHoleDiameter);
                    rotate([0,0,180]) translate([0,CaseWidth/2,0]) MountHolder(MountHolderThickness,MountHolderHoleDiameter);
                    }
                    else
                    {translate([0,CaseWidth/2,0]) MountHolder(MountHolderThickness,MountHolderHoleDiameter);}
                }
                if (ShowDeviceHolder1)
                {
                    translate([ DeviceHolder_X_Distance1/2+Offset_X_1, DeviceHolder_y_Distance1/2+Offset_Y_1,-0.01]) DeviceHolder("LightSalmon",ScrewCylinderHeight1,ScrewCylinderDiameter1,ScrewHoleDiameter1,DeviceHolder1_UseThread,DeviceHolder1_ScrewSize,DeviceHolder1_ThreadPitch,ThreadAngle,ThreadFitComp);
                    translate([-DeviceHolder_X_Distance1/2+Offset_X_1,-DeviceHolder_y_Distance1/2+Offset_Y_1,-0.01]) DeviceHolder("LightSalmon",ScrewCylinderHeight1,ScrewCylinderDiameter1,ScrewHoleDiameter1,DeviceHolder1_UseThread,DeviceHolder1_ScrewSize,DeviceHolder1_ThreadPitch,ThreadAngle,ThreadFitComp);
                    translate([ DeviceHolder_X_Distance1/2+Offset_X_1,-DeviceHolder_y_Distance1/2+Offset_Y_1,-0.01]) DeviceHolder("LightSalmon",ScrewCylinderHeight1,ScrewCylinderDiameter1,ScrewHoleDiameter1,DeviceHolder1_UseThread,DeviceHolder1_ScrewSize,DeviceHolder1_ThreadPitch,ThreadAngle,ThreadFitComp);
                    translate([-DeviceHolder_X_Distance1/2+Offset_X_1, DeviceHolder_y_Distance1/2+Offset_Y_1,-0.01]) DeviceHolder("LightSalmon",ScrewCylinderHeight1,ScrewCylinderDiameter1,ScrewHoleDiameter1,DeviceHolder1_UseThread,DeviceHolder1_ScrewSize,DeviceHolder1_ThreadPitch,ThreadAngle,ThreadFitComp);
                }
                if (ShowDeviceHolder2)
                {
                    translate([ DeviceHolder_X_Distance2/2+Offset_X_2, DeviceHolder_y_Distance2/2+Offset_Y_2,-0.01]) DeviceHolder("Khaki",ScrewCylinderHeight2,ScrewCylinderDiameter2,ScrewHoleDiameter2);
                    translate([-DeviceHolder_X_Distance2/2+Offset_X_2,-DeviceHolder_y_Distance2/2+Offset_Y_2,-0.01]) DeviceHolder("Khaki",ScrewCylinderHeight2,ScrewCylinderDiameter2,ScrewHoleDiameter2);
                    translate([ DeviceHolder_X_Distance2/2+Offset_X_2,-DeviceHolder_y_Distance2/2+Offset_Y_2,-0.01]) DeviceHolder("Khaki",ScrewCylinderHeight2,ScrewCylinderDiameter2,ScrewHoleDiameter2);
                    translate([-DeviceHolder_X_Distance2/2+Offset_X_2, DeviceHolder_y_Distance2/2+Offset_Y_2,-0.01]) DeviceHolder("Khaki",ScrewCylinderHeight2,ScrewCylinderDiameter2,ScrewHoleDiameter2);
                }
                if (ShowDeviceHolder3)
                {
                    translate([ DeviceHolder_X_Distance3/2+Offset_X_3, DeviceHolder_y_Distance3/2+Offset_Y_3,-0.01]) DeviceHolder("PaleGreen",ScrewCylinderHeight3,ScrewCylinderDiameter3,ScrewHoleDiameter3);
                    translate([-DeviceHolder_X_Distance3/2+Offset_X_3,-DeviceHolder_y_Distance3/2+Offset_Y_3,-0.01]) DeviceHolder("PaleGreen",ScrewCylinderHeight3,ScrewCylinderDiameter3,ScrewHoleDiameter3);
                    translate([ DeviceHolder_X_Distance3/2+Offset_X_3,-DeviceHolder_y_Distance3/2+Offset_Y_3,-0.01]) DeviceHolder("PaleGreen",ScrewCylinderHeight3,ScrewCylinderDiameter3,ScrewHoleDiameter3);
                    translate([-DeviceHolder_X_Distance3/2+Offset_X_3, DeviceHolder_y_Distance3/2+Offset_Y_3,-0.01]) DeviceHolder("PaleGreen",ScrewCylinderHeight3,ScrewCylinderDiameter3,ScrewHoleDiameter3);
                }

                // --- BATTERY HOLDER MOUNTING BOSSES (printed M3 threads, threads.scad) ---
                // Solid bosses rising from the inner floor. Thread is cut into these (as a blind hole,
                // see the difference() below) instead of relying on the battery holder's thin 2mm tab
                // or a separate nut. Floor stays fully sealed on the exterior underside.
                // V7.12 PhuNguyenPT: dropped the hard-coded $fn=30 on these two cylinders - that made
                // this boss visibly more faceted/polygonal than DeviceHolder()'s boss cylinder (the one
                // PCB1 rests on), which has no $fn override and so inherits the global $fn=80 set near
                // the top of the file. Removing the override here makes both bosses use the same
                // circumference resolution.
                color("SteelBlue") for (i = [0 : 2]) {
                    translate([BatteryHolder_LeftX + BatteryHolderHoleGapFromWidthEdge, BatteryHolder_BottomY + BatteryHolderHoleGapFromLengthEdge + (i * BatteryHolderHoleSpacing), BottomTopThickness])
                        cylinder(h=BatteryHolderBossHeight, d=BatteryHolderBossDiameter, center=false);

                    translate([BatteryHolder_RightX - BatteryHolderHoleGapFromWidthEdge, BatteryHolder_BottomY + BatteryHolderHoleGapFromLengthEdge + (i * BatteryHolderHoleSpacing), BottomTopThickness])
                        cylinder(h=BatteryHolderBossHeight, d=BatteryHolderBossDiameter, center=false);
                }
            }

            // Add here cut outs
            // --- BATTERY HOLDER MOUNTING HOLES: printed M3 threads (threads.scad), no nut needed ---
            // Blind thread cut into each boss added in the union() above. Starts just above the floor
            // (leaves the floor's exterior underside solid/sealed) and reaches to the top of the boss
            // so the screw can enter from inside the case, pass as clearance through the battery
            // holder's thin ~2mm tab, and thread straight into the boss.
            for (i = [0 : 2]) {
                // "Left" row (Nearest to the bottom-left case corner)
                translate([BatteryHolder_LeftX + BatteryHolderHoleGapFromWidthEdge, BatteryHolder_BottomY + BatteryHolderHoleGapFromLengthEdge + (i * BatteryHolderHoleSpacing), BottomTopThickness+0.01])
                    ScrewThread(1.01*MetricScrewSize+1.25*ThreadFitComp, BatteryHolderBossHeight, ThreadPitch, ThreadAngle, ThreadFitComp);

                // "Right" row (Separated by the longest distance / length)
                translate([BatteryHolder_RightX - BatteryHolderHoleGapFromWidthEdge, BatteryHolder_BottomY + BatteryHolderHoleGapFromLengthEdge + (i * BatteryHolderHoleSpacing), BottomTopThickness+0.01])
                    ScrewThread(1.01*MetricScrewSize+1.25*ThreadFitComp, BatteryHolderBossHeight, ThreadPitch, ThreadAngle, ThreadFitComp);
            }
            // -------------------------------------------------

            if (ShowBottomMountingHoles)
            {
                translate([BottomMountingHoles_X_Distance/2+BottomMountingHolesOffset_X,BottomMountingHoles_y_Distance/2+BottomMountingHolesOffset_Y,-0.01])cylinder(h=BottomTopThickness+0.02,d=BottomMountingHolesDiameter,center = false);
                translate([BottomMountingHoles_X_Distance/2+BottomMountingHolesOffset_X,-BottomMountingHoles_y_Distance/2+BottomMountingHolesOffset_Y,-0.01])cylinder(h=BottomTopThickness+0.02,d=BottomMountingHolesDiameter,center = false);
                translate([-BottomMountingHoles_X_Distance/2+BottomMountingHolesOffset_X,BottomMountingHoles_y_Distance/2+BottomMountingHolesOffset_Y,-0.01])cylinder(h=BottomTopThickness+0.02,d=BottomMountingHolesDiameter,center = false);
                translate([-BottomMountingHoles_X_Distance/2+BottomMountingHolesOffset_X,-BottomMountingHoles_y_Distance/2+BottomMountingHolesOffset_Y,-0.01])cylinder(h=BottomTopThickness+0.02,d=BottomMountingHolesDiameter,center = false);
            }

            if (SeeNutCut)           { color("red") translate([0,0,CaseHeight/2+CaseHeight-CutFromTop-NutSink]) cube([Caselength+0.1,CaseWidth+0.1,CaseHeight],center=true);}
            if (SeeGrooveRidgeScrew)
            {
                color("red") translate([Caselength-CaseRoundingRadius,0,(CaseHeight+0.1)/2-0.05])   cube([Caselength+0.1,CaseWidth*2+0.1,CaseHeight+0.1],center=true);
                color("red") translate([0,+CaseWidth/2-CaseWidth*2+CaseRoundingRadius,(CaseHeight+0.1)/2-0.05])   cube([Caselength+0.1,CaseWidth*2+0.1,CaseHeight+0.1],center=true);
            }

            if (ShowSideWallHoles_A)
            {
                if (CountOfSideWallHoles_A==1)
                {    translate([Caselength/2-SideWallThickness/2,SideWallHolePosition_A,SideWallHoleOffset_Z_A+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_A,center = true);}
                if (CountOfSideWallHoles_A==2)
                {
                    translate([Caselength/2-SideWallThickness/2, SideWallHoleDistance_A/2+SideWallHolePosition_A,SideWallHoleOffset_Z_A+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_A,center = true);
                    translate([Caselength/2-SideWallThickness/2,-SideWallHoleDistance_A/2+SideWallHolePosition_A,SideWallHoleOffset_Z_A+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_A,center = true);
                }
            if (CountOfSideWallHoles_A==3)
            {
                    // 1. Top PG9 Hole: Shifted inward by its radius PLUS an extra 8mm to clear the corner
                    translate([Caselength/2-SideWallThickness/2, BreadboardOffset_Y + BreadboardWidth/2 - (15.2 + 2*ClearanceGap)/2 - 8, SideWallHoleOffset_Z_A+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=15.2 + 2*ClearanceGap,center = true);

                    // 2. PG7 Hole: Aligned with PCB1 (solar controller)
                    if (!XAdditionalScrew) {
                        translate([Caselength/2-SideWallThickness/2, Offset_Y_1+10, SideWallHoleOffset_Z_A+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_A,center = true);
                    }

                    // 3. Bottom PG9 Hole: Shifted inward by its radius PLUS an extra 8mm to clear the corner
                    translate([Caselength/2-SideWallThickness/2, BreadboardOffset_Y - BreadboardWidth/2 + (15.2 + 2*ClearanceGap)/2 + 8, SideWallHoleOffset_Z_A+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=15.2 + 2*ClearanceGap,center = true);
                }
            }

            if (ShowSideWallHoles_B)
            {
                if (CountOfSideWallHoles_B==1)
                {    translate([SideWallHolePosition_B,-CaseWidth/2+SideWallThickness/2,SideWallHoleOffset_Z_B+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_B,center = true);}
                if (CountOfSideWallHoles_B==2)
                {
                    translate([SideWallHoleDistance_B/2+SideWallHolePosition_B,-CaseWidth/2+SideWallThickness/2, SideWallHoleOffset_Z_B+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_B,center = true);
                    translate([-SideWallHoleDistance_B/2+SideWallHolePosition_B,-CaseWidth/2+SideWallThickness/2, SideWallHoleOffset_Z_B+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_B,center = true);
                }
                if (CountOfSideWallHoles_B==3)
                {
                    translate([SideWallHolePosition_B,-CaseWidth/2+SideWallThickness/2, SideWallHoleOffset_Z_B+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_B,center = true);
                    translate([SideWallHoleDistance_B+SideWallHolePosition_B,-CaseWidth/2+SideWallThickness/2, SideWallHoleOffset_Z_B+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_B,center = true);
                    translate([-SideWallHoleDistance_B+SideWallHolePosition_B,-CaseWidth/2+SideWallThickness/2,SideWallHoleOffset_Z_B+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_B,center = true);
                }
            }

            if (ShowSideWallHoles_C)
            {
                if (CountOfSideWallHoles_C==1)
                {    translate([-Caselength/2+SideWallThickness/2,SideWallHolePosition_C,SideWallHoleOffset_Z_C+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_C,center = true);}
                if (CountOfSideWallHoles_C==2)
                {
                    translate([-Caselength/2+SideWallThickness/2, SideWallHoleDistance_C/2+SideWallHolePosition_C,SideWallHoleOffset_Z_C+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_C,center = true);
                    translate([-Caselength/2+SideWallThickness/2,-SideWallHoleDistance_C/2+SideWallHolePosition_C,SideWallHoleOffset_Z_C+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_C,center = true);
                }
                if (CountOfSideWallHoles_C==3)
                {
                    translate([-Caselength/2+SideWallThickness/2, SideWallHolePosition_C,SideWallHoleOffset_Z_C+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_C,center = true);
                    translate([-Caselength/2+SideWallThickness/2, SideWallHoleDistance_C+SideWallHolePosition_C,SideWallHoleOffset_Z_C+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_C,center = true);
                    translate([-Caselength/2+SideWallThickness/2,-SideWallHoleDistance_C+SideWallHolePosition_C,SideWallHoleOffset_Z_C+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,90]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_C,center = true);
                }
            }

            if (ShowSideWallHoles_D)
            {
                if (CountOfSideWallHoles_D==1)
                {    translate([SideWallHolePosition_D,CaseWidth/2-SideWallThickness/2,SideWallHoleOffset_Z_D+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_D,center = true);}
                if (CountOfSideWallHoles_D==2)
                {
                    translate([SideWallHoleDistance_D/2+SideWallHolePosition_D,CaseWidth/2-SideWallThickness/2, SideWallHoleOffset_Z_D+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_D,center = true);
                    translate([-SideWallHoleDistance_D/2+SideWallHolePosition_D,CaseWidth/2-SideWallThickness/2, SideWallHoleOffset_Z_D+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_D,center = true);
                }
                if (CountOfSideWallHoles_D==3)
                {
                    translate([SideWallHolePosition_D,CaseWidth/2-SideWallThickness/2, SideWallHoleOffset_Z_D+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_D,center = true);
                    translate([SideWallHoleDistance_D+SideWallHolePosition_D,CaseWidth/2-SideWallThickness/2, SideWallHoleOffset_Z_D+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_D,center = true);
                    translate([-SideWallHoleDistance_D+SideWallHolePosition_D,CaseWidth/2-SideWallThickness/2,SideWallHoleOffset_Z_D+ BottomTopThickness+(CaseHeight-CutFromTop-BottomTopThickness)/2 ]) rotate([90,0,0]) cylinder(h=SideWallThickness+0.04,d=SideWallHoleDiameter_D,center = true);
                }
            }
        }
    }
}

module BodyTop () {
    if (ShowTop)
    {
        difference(){
            union(){
                rotate([  0,  0,  0]) BodyQuarterTop(Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness);
                rotate([  0,  0,180]) BodyQuarterTop(Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness);
                mirror([  0,  1,  0]) BodyQuarterTop(Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness);
                mirror([  1,  0  ,0]) BodyQuarterTop(Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness);

                // RJ45 friction-fit collars - added material for a snug grip along the module's inserted depth
                if (ShowEthernetPortCollar && ShowEthernetPort1) { color("SlateGray") EthernetPortCollar(EthernetPortWidth,EthernetPortHeight,EthernetPortInsideDepth,EthernetPortCollarWall,EthernetPort1_Wall,EthernetPort1_Position, CaseHeight-(EthernetPortDefaultZ+EthernetPort1_OffsetZ)); }
                if (ShowEthernetPortCollar && ShowEthernetPort2) { color("SlateGray") EthernetPortCollar(EthernetPortWidth,EthernetPortHeight,EthernetPortInsideDepth,EthernetPortCollarWall,EthernetPort2_Wall,EthernetPort2_Position, CaseHeight-(EthernetPortDefaultZ+EthernetPort2_OffsetZ)); }
            }
            if (SeeGrooveRidgeScrew)
            {
                color("red") translate([-Caselength+CaseRoundingRadius,0,(CaseHeight+0.1)/2-0.05]) cube([Caselength+0.1,CaseWidth+0.1,CaseHeight+0.1],center=true);
                color("red") translate([0,-CaseWidth+CaseRoundingRadius,(CaseHeight+0.1)/2-0.05]) cube([Caselength+0.1,CaseWidth+0.1,CaseHeight+0.1],center=true);
            }

            // RJ45 ethernet port cutouts - positioned in the top piece so they sit above the breadboard
            if (ShowEthernetPort1) { EthernetPortCut(EthernetPortWidth,EthernetPortHeight,EthernetPortInsideDepth,EthernetPort1_Wall,EthernetPort1_Position, CaseHeight-(EthernetPortDefaultZ+EthernetPort1_OffsetZ)); }
            if (ShowEthernetPort2) { EthernetPortCut(EthernetPortWidth,EthernetPortHeight,EthernetPortInsideDepth,EthernetPort2_Wall,EthernetPort2_Position, CaseHeight-(EthernetPortDefaultZ+EthernetPort2_OffsetZ)); }
        }
    }
}

module MountHolder (Thick,Hole) {
    translate([0,0,0.005]) difference(){

        if (MountHolderStyle==1){
            $fn=40;
           roundedBox([Caselength, MountHolderLenght*2, Thick*2], Thick/3, 0);
        }
        if (MountHolderStyle==2){
            $fn=60;
           roundedBox([Caselength, MountHolderLenght*2, Thick*2],CaseRoundingRadius , 1);
        }
        if (MountHolderStyle==3){
            roundedBox([Caselength, MountHolderLenght*2, Thick*2],0 , 2);
        }
        if((MountHolderStyle>0)&&(MountHolderStyle<4)){
            translate([0,0,-Thick/2-0.02]) cube([Caselength+0.02,MountHolderLenght*2+0.02,Thick+0.04],center=true);
            translate([0,-MountHolderLenght/2-CaseRoundingRadius,Thick/2+0.02]) cube([Caselength+0.02,MountHolderLenght+0.02,Thick+0.08],center=true);
            translate([0,-MountHolderLenght/2,Thick/2+0.02]) cube([Caselength-2*CaseRoundingRadius+0.02,MountHolderLenght+0.02,Thick+0.08],center=true);
            translate([0,-CaseRoundingRadius,MountHolderThickness/2-0.02]) translate(ScrewAddYPos) cylinder(h=MountHolderThickness+0.06,d=ScrewHoleDia,center = true);
            mirror([  1,  0,  0]) translate([0,-CaseRoundingRadius,MountHolderThickness/2-0.02]) translate(ScrewAddYPos) cylinder(h=MountHolderThickness+0.06,d=ScrewHoleDia,center = true);

            if (CountOfMountHolderHoles>1){
                translate([Caselength/2-Hole-Thick/3,Hole*1.5,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
                translate([-Caselength/2+Hole+Thick/3,Hole*1.5,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
            }
            if (CountOfMountHolderHoles!=2){
                translate([0,Hole*1.5,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
            }
        }
    }
    if (MountHolderStyle==4){
        HolderRad=Hole/2;
        HolderWidth=4*Hole;
        translate([0,MountHolderLenght,0]) difference(){
            union(){
                translate([0,-MountHolderLenght+HolderRad,0]) roundedBox([HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
                translate([-HolderWidth/2+HolderRad,0,0]) rotate([0,0,-45]) translate([HolderWidth-HolderRad,-MountHolderLenght+HolderRad,0]) roundedBox([2*HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
                translate([ HolderWidth/2-HolderRad,0,0]) rotate([0,0,45]) translate([-HolderWidth+HolderRad,-MountHolderLenght+HolderRad,0]) roundedBox([2*HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
            }
            // FIX 1: Shifted Y-cut by -0.1mm to create clean overlap into the case wall
            translate([0,-(3*HolderWidth)/2-MountHolderLenght - 0.1,-0.02]) cube([10*HolderWidth,3*HolderWidth,Thick*4+0.06],center=true);
            // FIX 2: Oversized the X/Y footprint and centered to catch ALL bottom geometry
            translate([0, 0, -Thick-0.02]) cube([10*HolderWidth,10*HolderWidth,Thick*2],center=true);

            translate([0,-MountHolderLenght+Hole*1.8,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
        }
    }
    if (MountHolderStyle==5){
        HolderRad=Hole/2;
        HolderWidth=4*Hole;
        translate([0,MountHolderLenght,0]) difference()
        {
            union(){
                $fn=40;
                translate([0,-MountHolderLenght+HolderRad,0]) roundedBox([HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
                translate([-HolderWidth/2+HolderRad,0,0]) rotate([0,0,-45]) translate([HolderWidth-HolderRad,-MountHolderLenght+HolderRad,0]) roundedBox([2*HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
                translate([ HolderWidth/2-HolderRad,0,0]) rotate([0,0,45]) translate([-HolderWidth+HolderRad,-MountHolderLenght+HolderRad,0]) roundedBox([2*HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
            }
            // FIX 1: Shifted Y-cut by -0.1mm to create clean overlap into the case wall
            translate([0,-(3*HolderWidth)/2-MountHolderLenght - 0.1,-0.02]) cube([10*HolderWidth,3*HolderWidth,Thick*4+0.06],center=true);
            // FIX 2: Oversized the X/Y footprint and centered to catch ALL bottom geometry
            translate([0, 0, -Thick-0.02]) cube([10*HolderWidth,10*HolderWidth,Thick*2],center=true);

            hull(){
                translate([-Hole/1.1,-MountHolderLenght+Hole*1.6,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
                translate([+Hole/1.1,-MountHolderLenght+Hole*1.6,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
            }
        }
        rotate([0,0,180]) translate([0,MountHolderLenght+CaseWidth,0]) difference()
        {
            union(){
                $fn=40;
                translate([0,-MountHolderLenght+HolderRad,0]) roundedBox([HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
                translate([-HolderWidth/2+HolderRad,0,0]) rotate([0,0,-45]) translate([HolderWidth-HolderRad,-MountHolderLenght+HolderRad,0]) roundedBox([2*HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
                translate([ HolderWidth/2-HolderRad,0,0]) rotate([0,0,45]) translate([-HolderWidth+HolderRad,-MountHolderLenght+HolderRad,0]) roundedBox([2*HolderWidth, MountHolderLenght*2, Thick*2],HolderRad , 0);
            }
            // FIX 1: Shifted Y-cut by -0.1mm to create clean overlap into the case wall
            translate([0,-(3*HolderWidth)/2-MountHolderLenght - 0.1,-0.02]) cube([10*HolderWidth,3*HolderWidth,Thick*4+0.06],center=true);
            // FIX 2: Oversized the X/Y footprint and centered to catch ALL bottom geometry
            translate([0, 0, -Thick-0.02]) cube([10*HolderWidth,10*HolderWidth,Thick*2],center=true);

            hull(){
                    translate([0,-MountHolderLenght+Hole*1.6+Hole/1.8,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
                    translate([0,-MountHolderLenght+Hole*1.6-Hole/1.8,MountHolderThickness/2-0.02]) cylinder(h=MountHolderThickness+0.06,d=Hole,center = true);
            }
        }
    }
}

module ShowSizes () {
    echo ();
    echo (str(" Stable and waterproof OpenSCAD case by pbtec V7."));
    echo ();
    echo (str(" Render quality : ",$fn));
    echo ();
    echo (str(" --> Case outer dimensions: "));
    echo (str(" Length : ",Caselength,"mm "));
    echo (str(" Width : ",CaseWidth,"mm "));
    echo (str(" High : ",CaseHeight,"mm "));
    echo (str(" Top (upper piece) high : ",CutFromTop,"mm "));
    echo (str(" Bottom (lower piece) high : ",CaseHeight-CutFromTop,"mm "));
    echo (str(" Side wall thickness : ",SideWallThickness,"mm "));
    echo (str(" Bottom & top wall thickness : ",BottomTopThickness,"mm "));
    echo (str(" Case rounding radius : ",CaseRoundingRadius,"mm "));
    echo ();
    echo (str(" --> Case inner dimensions : "));
    echo (str(" X : Wall to wall : ",Caselength-2*SideWallThickness,"mm "));
    echo (str(" X : Screw cylinder to screw cylinder : ",Caselength-4*CaseRoundingRadius,"mm "));
    echo (str(" Y : Wall to wall : ",CaseWidth-2*SideWallThickness,"mm "));
    echo (str(" Y : Screw cylinder to screw cylinder : ",CaseWidth-4*CaseRoundingRadius,"mm "));
    echo (str(" Top to bottom  : ",CaseHeight-2*BottomTopThickness,"mm "));
    echo ();
    echo (str(" --> Breadboard fit check (",BreadboardLength,"x",BreadboardWidth,"x",BreadboardHeight,"mm): "));
    echo (str(" X clearance (Length) : ",(Caselength-2*SideWallThickness)-BreadboardLength,"mm "));
    echo (str(" Y clearance (Width)  : ",(CaseWidth-2*SideWallThickness)-BreadboardWidth,"mm "));
    echo (str(" Z clearance (Height) : ",(CaseHeight-2*BottomTopThickness)-BreadboardHeight,"mm "));
    echo ();
    echo (str(" --> Battery holder fit check (",BatteryHolderLength,"x",BatteryHolderWidth,"mm footprint, ",BatteryHolderFloorThickness+max(BatteryHolderWidthWallHeight,BatteryHolderLengthWallHeight),"mm tall - floor ",BatteryHolderFloorThickness,"mm + tallest wall ",max(BatteryHolderWidthWallHeight,BatteryHolderLengthWallHeight),"mm): "));
    echo (str(" X clearance (Length) : ",(Caselength-2*SideWallThickness)-BatteryHolderLength,"mm "));
    echo (str(" Y clearance (Width)  : ",(CaseWidth-2*SideWallThickness)-BatteryHolderWidth,"mm "));
    echo (str(" Z clearance (Height), holder resting on top of the ",BatteryHolderBossHeight,"mm mounting boss : ",(CaseHeight-2*BottomTopThickness-BatteryHolderBossHeight)-(BatteryHolderFloorThickness+max(BatteryHolderWidthWallHeight,BatteryHolderLengthWallHeight)),"mm "));
    echo (str(" Battery slot width (each of 3 slots, between ",BatteryHolderSlotWallThickness,"mm dividers) : ",(BatteryHolderWidth-2*BatteryHolderWidthWallThickness-2*BatteryHolderSlotWallThickness)/3,"mm "));
    echo (str(" Gap between breadboard and battery holder (Y) : ",(BreadboardOffset_Y-BreadboardWidth/2)-(BatteryHolderOffset_Y+BatteryHolderWidth/2),"mm (must be > 0, no overlap)"));
    echo ();
    echo (str(" --> RJ45 ethernet port cutouts (",EthernetPortWidth,"x",EthernetPortHeight,"mm opening, ",EthernetPortInsideDepth,"mm inside clearance): "));
    echo (str(" Cut into the TOP piece - centered height (global Z) : ",EthernetPortDefaultZ,"mm "));
    echo (str(" Breadboard top (global Z) : ",BottomTopThickness+BreadboardHeight,"mm "));
    echo (str(" Port 1 : wall ",EthernetPort1_Wall," enabled=",ShowEthernetPort1));
    echo (str(" Port 2 : wall ",EthernetPort2_Wall," enabled=",ShowEthernetPort2));
    echo (str(" Friction-fit collar : enabled=",ShowEthernetPortCollar," wall thickness=",EthernetPortCollarWall,"mm "));
    echo ();
    echo (str(" Screw dimensions : "));
    echo (str(" Metric Screw size: m",MetricScrewSize));
    echo (str(" Screw hole diameter : ",ScrewHoleDia,"mm "));
    echo (str(" Screw head diameter : ",ScrewHeadDia,"mm "));
    echo (str(" X : Additional screw (3rd)) : ",XAdditionalScrew));
    echo (str(" Y : Additional screws (3rd) : ",YAdditionalScrew));
    echo (str(" --> Check if you have screws within the following size : "));

    if(NutStyle ==1)  // Hot melt copper nuts
    {
        if (ScrewHeadType == 1 ) // countersunk screws and flat head screws without counter
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=HolelengthHotMeltNut;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Countersunk Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Countersunk Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
        if (ScrewHeadType== 2)
        {
            vTopMin=CutFromTop-ScrewHeadHeight;
            vTopMax=CutFromTop-ScrewHeadHeight;
            vBottomMin=HolelengthHotMeltNut;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }

        if (ScrewHeadType == 3) // countersunk screws and flat head screws without counter
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=HolelengthHotMeltNut;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
    }

    if(NutStyle ==2) // standard nuts and square nuts
    {
        if (ScrewHeadType == 1) // countersunk screws and flat head screws without counter
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=NutSink+NutHigh;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Countersunk Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Countersunk Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
        if (ScrewHeadType== 2)
        {
            vTopMin=CutFromTop-ScrewHeadHeight;
            vTopMax=CutFromTop-ScrewHeadHeight;
            vBottomMin=NutSink+NutHigh;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }

        if (ScrewHeadType == 3) // countersunk screws and flat head screws without counter
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=NutSink+NutHigh;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
    }

    if(NutStyle == 3 || NutStyle == 4) // standard nuts and square nuts
    {
        if (ScrewHeadType == 1) // countersunk screws and flat head screws without counter
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=NutSink+SquareNutHigh;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Countersunk Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Countersunk Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
        if (ScrewHeadType== 2)
        {
            vTopMin=CutFromTop-ScrewHeadHeight;
            vTopMax=CutFromTop-ScrewHeadHeight;
            vBottomMin=NutSink+SquareNutHigh;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
        if (ScrewHeadType == 3) // cflat head screws without counter
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=NutSink+SquareNutHigh;
            vBottomMax=HoleDeepness;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }

    }

    if(NutStyle == 5) // Printed internal M3 threads (threads.scad), no insert/nut - screw threads directly into the corner boss
    {
        echo (str(" --> Corner bosses use PRINTED INTERNAL THREADS (threads.scad):"));
        echo (str(" Thread pitch : ",ThreadPitch,"mm  |  Thread tooth angle : ",ThreadAngle,"deg  |  Fit tolerance : ",ThreadFitComp,"mm"));
        echo (str(" Effective internal thread cut diameter : ",1.01*MetricScrewSize+1.25*ThreadFitComp,"mm  (nominal m",MetricScrewSize,", threads.scad oversize formula)"));
        echo (str(" Thread engagement length : ",ThreadEngagementLength,"mm"));
        if (ScrewHeadType == 1) // countersunk / flat head screws
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=2*MetricScrewSize; // conservative min thread engagement, ~2x nominal diameter for FDM plastic
            vBottomMax=ThreadEngagementLength;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Countersunk (flat head) Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Countersunk (flat head) Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
        if (ScrewHeadType == 2)
        {
            vTopMin=CutFromTop-ScrewHeadHeight;
            vTopMax=CutFromTop-ScrewHeadHeight;
            vBottomMin=2*MetricScrewSize;
            vBottomMax=ThreadEngagementLength;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Round head or hex Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
        if (ScrewHeadType == 3)
        {
            vTopMin=CutFromTop;
            vTopMax=CutFromTop;
            vBottomMin=2*MetricScrewSize;
            vBottomMax=ThreadEngagementLength;
            vMin=vTopMin+vBottomMin;
            vMax=vTopMax+vBottomMax;
            echo (str(" --> Exposed head Screw m",MetricScrewSize , " max length : ",vMax, "mm"));
            echo (str(" --> Exposed head Screw m",MetricScrewSize , " min length : ",vMin, "mm"));
        }
    }
}

module GrooveStraight (length) {
   color("orange") translate([length/2,0,-(GrooveDepth+AddGrooveDepthForSealing)/2]) cube([length,GrooveWidth,GrooveDepth+AddGrooveDepthForSealing],center=true);
}

module GrooveCurved (Angle,Rad) {
    color("orange") difference(){
        translate([0,0,-(GrooveDepth+AddGrooveDepthForSealing)]) pie(Rad+(GrooveWidth)/2, Angle, GrooveDepth+AddGrooveDepthForSealing, spin=0);
        translate([-0.01,-0.01,-(GrooveDepth+AddGrooveDepthForSealing+0.02)]) pie(Rad-(GrooveWidth)/2, Angle, GrooveDepth+AddGrooveDepthForSealing+0.04, spin=0);
    }
}

module RidgeStraight (length) {
    color("orange") translate([length/2,0,(RidgeHeight)/2]) cube([length,RidgeWidth,RidgeHeight],center=true);
}

module RidgeCurved (Angle,Rad) {
    color("orange") difference(){
        translate([0,0,0])         pie(Rad+(RidgeWidth)/2, Angle, RidgeHeight, spin=0);
        translate([-0.01,-0.01,-0.02]) pie(Rad-(RidgeWidth)/2, Angle, RidgeHeight+0.04, spin=0);
    }
}

module GasketRidgeStraight (length) {
    color("DarkSlateGray")
    translate([length/2,0,RidgeHeight+(GasketHeight)/2])
        cube([length,GasketRidgeWidth,GasketHeight],center=true);
}

module GasketRidgeCurved (Angle,Rad) {
    color("DarkSlateGray") difference(){
        translate([0,0,RidgeHeight])
            pie(Rad+(GasketRidgeWidth)/2, Angle, GasketHeight, spin=0);
        translate([-0.01,-0.01,RidgeHeight-0.02])
            pie(Rad-(GasketRidgeWidth)/2, Angle, GasketHeight+0.04, spin=0);
    }
}

module Gasket () {
    rotate([  0,  0,  0]) GasketQuarter();
    rotate([  0,  0,180]) GasketQuarter();
    mirror([  0,  1,  0]) GasketQuarter();
    mirror([  1,  0  ,0]) GasketQuarter();
}

module GasketQuarter () {
            translate([CaseRoundingRadius+ScrewHoleDia/2-0.01,CaseWidth/2-OuterBorder-GrooveWidth/2-0.01,0])   GasketRidgeStraight(Caselength/2-3*CaseRoundingRadius-ScrewHoleDia+0.03);
            translate([Caselength/2-OuterBorder-GrooveWidth/2-0.01,CaseRoundingRadius+ScrewHoleDia/2-0.02,0]) translate([0,0,0]) rotate([0,0,90]) GasketRidgeStraight(CaseWidth/2-3*CaseRoundingRadius-ScrewHoleDia+0.04);
            translate(ScrewCornerPos) rotate([0,0,180]) GasketRidgeCurved(90,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
            translate([-ScrewHoleDia-SideWallThickness+0.02,-0.01,0]) translate(ScrewCornerPos)GasketRidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            translate([-0.01,-ScrewHoleDia-SideWallThickness+0.00,0]) translate(ScrewCornerPos)GasketRidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            if (XAdditionalScrew){
                translate(ScrewAddXPos)   rotate([0,0,180]) GasketRidgeCurved(180,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
                translate([-ScrewHoleDia-SideWallThickness+0.04,-0.01,0]) translate(ScrewAddXPos) GasketRidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
                translate([ScrewHoleDia+SideWallThickness,-0.01,0]) translate(ScrewAddXPos) rotate([0,0,90]) GasketRidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            }
            else { translate([-0.01,CaseWidth/2-OuterBorder-GrooveWidth/2-0.01,0])  GasketRidgeStraight(Caselength/2-2*CaseRoundingRadius-ScrewHoleDia/2+0.05);}
            if (YAdditionalScrew){
                translate(ScrewAddYPos)   rotate([0,0,90])  GasketRidgeCurved(180,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
                translate([-0.01,-ScrewHoleDia-SideWallThickness-0.01,0]) translate(ScrewAddYPos) GasketRidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
                translate([-0.01,ScrewHoleDia+SideWallThickness-0.01,0]) translate(ScrewAddYPos) rotate([0,0,270]) GasketRidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            }
            else { translate([Caselength/2-OuterBorder-GrooveWidth/2-0.01,-0.01,0])rotate([0,0,90]) GasketRidgeStraight(CaseWidth/2-2*CaseRoundingRadius-ScrewHoleDia/2+0.05);}
}


module BodyQuarterBottom (Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness) {
    difference(){
        union(){
            color("SteelBlue")BodyQuarter(Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness);

            translate([CaseRoundingRadius+ScrewHoleDia/2-0.01,CaseWidth/2-OuterBorder-GrooveWidth/2-0.01,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,0])  RidgeStraight(Caselength/2-3*CaseRoundingRadius-ScrewHoleDia+0.03);
            translate([Caselength/2-OuterBorder-GrooveWidth/2-0.01,CaseRoundingRadius+ScrewHoleDia/2-0.02,CutFromTop+0.01]) translate([0,0,0]) rotate([0,0,90]) RidgeStraight(CaseWidth/2-3*CaseRoundingRadius-ScrewHoleDia+0.04);
            translate([0,0,CutFromTop+0.01]) translate(ScrewCornerPos) rotate([0,0,180]) RidgeCurved(90,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
            translate([-ScrewHoleDia-SideWallThickness+0.02,-0.01,CutFromTop+0.01]) translate(ScrewCornerPos) rotate([0,0,0]) RidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            translate([-0.01,-ScrewHoleDia-SideWallThickness+0.00,CutFromTop+0.01]) translate(ScrewCornerPos) rotate([0,0,0]) RidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            if (XAdditionalScrew){
                translate([0,0,CutFromTop+0.01]) translate(ScrewAddXPos)   rotate([0,0,180])                                RidgeCurved(180,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
                translate([-ScrewHoleDia-SideWallThickness+0.04,-0.01,CutFromTop+0.01]) translate(ScrewAddXPos) rotate([0,0,0]) RidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
                translate([ScrewHoleDia+SideWallThickness,-0.01,CutFromTop+0.01]) translate(ScrewAddXPos) rotate([0,0,90])     RidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            }
            else { translate([-0.01,CaseWidth/2-OuterBorder-GrooveWidth/2-0.01,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,0]) RidgeStraight(Caselength/2-2*CaseRoundingRadius-ScrewHoleDia/2+0.05);}
            if (YAdditionalScrew){
                translate([0,0,CutFromTop+0.01]) translate(ScrewAddYPos)   rotate([0,0,90])  RidgeCurved(180,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
                translate([-0.01,-ScrewHoleDia-SideWallThickness-0.01,CutFromTop+0.01]) translate(ScrewAddYPos) rotate([0,0,0]) RidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
                translate([-0.01,ScrewHoleDia+SideWallThickness-0.01,CutFromTop+0.01]) translate(ScrewAddYPos) rotate([0,0,270]) RidgeCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            }
            else { translate([Caselength/2-OuterBorder-GrooveWidth/2-0.01,-0.01,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,90]) RidgeStraight(CaseWidth/2-2*CaseRoundingRadius-ScrewHoleDia/2+0.05);}
        }
        if(NutStyle == 1)
        {translate(ScrewCornerPos) translate([0,0,CutFromTop-HoleDeepness+0.01]) cylinder(h=HoleDeepness ,d=HoleDiaThread,center = false);}
        if(NutStyle == 2) {translate(ScrewCornerPos) translate([0,0,CutFromTop+0.01]) NutCut(CutFromTop,NutHigh,NutDia);}
        if(NutStyle == 3) {translate(ScrewCornerPos) translate([0,0,CutFromTop+0.01]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,1);}
        if(NutStyle == 4)
        {
            if (EdgeSquareNutInsertFrom_X) {translate(ScrewCornerPos) translate([0,0,CutFromTop+0.01]) rotate([0,0, 0]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,0);}
            else                           {translate(ScrewCornerPos) translate([0,0,CutFromTop+0.01]) rotate([0,0,90]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,0);}
        }
        if(NutStyle == 5)
        {
            // Printed internal M3 thread cut (threads.scad ScrewThread()), corner boss only.
            // Positioned exactly like the NutStyle 1 plain hole: engagement zone
            // sits just below the split line (top face of the bottom piece),
            // so the flat-head screw coming down through the top piece threads
            // straight into it.
            // Oversize formula (1.01*outer_diam + 1.25*tolerance) matches exactly what
            // threads.scad's own ScrewHole() wrapper applies - reproduced directly here so
            // the cut composes with this file's existing single difference() block.
            translate(ScrewCornerPos) translate([0,0,CutFromTop-ThreadEngagementLength+0.01])
                ScrewThread(1.01*MetricScrewSize+1.25*ThreadFitComp, ThreadEngagementLength, ThreadPitch, ThreadAngle, ThreadFitComp);
        }

        if (XAdditionalScrew){
            if (NutStyle ==1){translate(ScrewAddXPos) translate([0,0,CutFromTop-HoleDeepness+0.01]) cylinder(h=HoleDeepness ,d=HoleDiaThread,center = false);}
            if (NutStyle ==2){translate(ScrewAddXPos) translate([0,0,CutFromTop+0.01]) NutCut(CutFromTop,NutHigh,NutDia);}
            if(NutStyle == 3) {translate(ScrewAddXPos) translate([0,0,CutFromTop+0.01]) rotate([0,0,90]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,1);}
            if(NutStyle == 4) {translate(ScrewAddXPos) translate([0,0,CutFromTop+0.01]) rotate([0,0,90]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,0);}
            if(NutStyle == 5) {translate(ScrewAddXPos) translate([0,0,CutFromTop-ThreadEngagementLength+0.01]) ScrewThread(1.01*MetricScrewSize+1.25*ThreadFitComp, ThreadEngagementLength, ThreadPitch, ThreadAngle, ThreadFitComp);}

        }
        if (YAdditionalScrew){
            if (NutStyle ==1){translate(ScrewAddYPos) translate([0,0,CutFromTop-HoleDeepness+0.01]) cylinder(h=HoleDeepness ,d=HoleDiaThread,center = false);}
            if (NutStyle ==2){translate(ScrewAddYPos) translate([0,0,CutFromTop+0.01]) NutCut(CutFromTop,NutHigh,NutDia);}
            if(NutStyle == 3) {translate(ScrewAddYPos) translate([0,0,CutFromTop+0.01]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,1);}
            if(NutStyle == 4) {translate(ScrewAddYPos) translate([0,0,CutFromTop+0.01]) SquareNutCut(CutFromTop,SquareNutHigh,SquareNutSize,0);}
            if(NutStyle == 5) {translate(ScrewAddYPos) translate([0,0,CutFromTop-ThreadEngagementLength+0.01]) ScrewThread(1.01*MetricScrewSize+1.25*ThreadFitComp, ThreadEngagementLength, ThreadPitch, ThreadAngle, ThreadFitComp);}

        }
    }
}

module BodyQuarterTop (Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness) {

    difference()
    {
        union(){
            color("DarkCyan")BodyQuarter(Caselength,CaseWidth,CutFromTop,CaseRoundingRadius,SideWallThickness);
        }
        if (ScrewHeadType == 1)
        {
            translate(ScrewCornerPos) ScrewCut(MetricScrewSize,CutFromTop+0.01,0);
            if (XAdditionalScrew){
                translate(ScrewAddXPos) ScrewCut(MetricScrewSize,CutFromTop+0.01,0);
            }
            if (YAdditionalScrew){
                translate(ScrewAddYPos) ScrewCut(MetricScrewSize,CutFromTop+0.01,0);
            }
        }
        if (ScrewHeadType == 2)
        {
            translate(ScrewCornerPos) translate([0,0,-0.01]) cylinder(h=CutFromTop+0.02,d=ScrewHoleDia ,center = false);
            translate(ScrewCornerPos) translate([0,0,-0.01]) cylinder(h=ScrewHeadHeight+0.02,d=ScrewHeadDiameter ,center = false);
            //ScrewCut(MetricScrewSize,CutFromTop+0.01,0);
            if (XAdditionalScrew){
                translate(ScrewAddXPos)  translate([0,0,-0.01]) cylinder(h=CutFromTop+0.02,d=ScrewHoleDia ,center = false);
                translate(ScrewAddXPos)  translate([0,0,-0.01]) cylinder(h=ScrewHeadHeight+0.02,d=ScrewHeadDiameter ,center = false);
            }
            if (YAdditionalScrew){
                translate(ScrewAddYPos)  translate([0,0,-0.01]) cylinder(h=CutFromTop+0.02,d=ScrewHoleDia ,center = false);
                translate(ScrewAddYPos)  translate([0,0,-0.01]) cylinder(h=ScrewHeadHeight+0.02,d=ScrewHeadDiameter ,center = false);
            }
        }
        if (ScrewHeadType == 3)
        {
            translate(ScrewCornerPos) translate([0,0,-0.01]) cylinder(h=CutFromTop+0.02,d=ScrewHoleDia ,center = false);
            //ScrewCut(MetricScrewSize,CutFromTop+0.01,0);
            if (XAdditionalScrew){
                translate(ScrewAddXPos)  translate([0,0,-0.01]) cylinder(h=CutFromTop+0.02,d=ScrewHoleDia ,center = false);
            }
            if (YAdditionalScrew){
                translate(ScrewAddYPos)  translate([0,0,-0.01]) cylinder(h=CutFromTop+0.02,d=ScrewHoleDia ,center = false);
            }
        }

        translate([CaseRoundingRadius+ScrewHoleDia/2-0.01,CaseWidth/2-OuterBorder-GrooveWidth/2-0.01,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,0]) GrooveStraight(Caselength/2-3*CaseRoundingRadius-ScrewHoleDia+0.03);
        translate([Caselength/2-OuterBorder-GrooveWidth/2-0.01,CaseRoundingRadius+ScrewHoleDia/2-0.02,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,90]) GrooveStraight(CaseWidth/2-3*CaseRoundingRadius-ScrewHoleDia+0.04);
        translate([0,0,CutFromTop+0.01]) translate(ScrewCornerPos) rotate([0,0,180]) GrooveCurved(90,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
        translate([-ScrewHoleDia-SideWallThickness+0.02,-0.01,CutFromTop+0.01]) translate(ScrewCornerPos) rotate([0,0,0]) GrooveCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
        translate([-0.01,-ScrewHoleDia-SideWallThickness+0.0,CutFromTop+0.01]) translate(ScrewCornerPos) rotate([0,0,0]) GrooveCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
        if (XAdditionalScrew){
            translate([0,0,CutFromTop+0.01]) translate(ScrewAddXPos)   rotate([0,0,180])                                GrooveCurved(180,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
            translate([-ScrewHoleDia-SideWallThickness+0.04,-0.01,CutFromTop+0.01]) translate(ScrewAddXPos) rotate([0,0,0]) GrooveCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            translate([ScrewHoleDia+SideWallThickness,-0.01,CutFromTop+0.01]) translate(ScrewAddXPos) rotate([0,0,90])     GrooveCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
        }
        else { translate([-0.01,CaseWidth/2-OuterBorder-GrooveWidth/2-0.01,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,0]) GrooveStraight(Caselength/2-2*CaseRoundingRadius-ScrewHoleDia/2+0.07); }

        if (YAdditionalScrew){
            translate([0,0,CutFromTop+0.01]) translate(ScrewAddYPos)   rotate([0,0,90])  GrooveCurved(180,ScrewHoleDia/2+OuterBorder+GrooveWidth/2);
            translate([-0.01,-ScrewHoleDia-SideWallThickness-0.01,CutFromTop+0.01]) translate(ScrewAddYPos) rotate([0,0,0]) GrooveCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
            translate([-0.01,ScrewHoleDia+SideWallThickness-0.01,CutFromTop+0.01]) translate(ScrewAddYPos) rotate([0,0,270]) GrooveCurved(90,ScrewHoleDia/2+InnerBorder+GrooveWidth/2);
        }
        else { translate([Caselength/2-OuterBorder-GrooveWidth/2-0.01,-0.01,CutFromTop+0.01])  translate([0,0,0]) rotate([0,0,90]) GrooveStraight(CaseWidth/2-2*CaseRoundingRadius-ScrewHoleDia/2+0.057); }
    }
}

module BodyQuarter (L,W,H,Rad,Rand){
    cube([L/2-Rad,W/2,BottomTopThickness],center=false); // Ground
    cube([L/2,W/2-Rad,BottomTopThickness],center=false); // Ground
    if (CaseRadius < CaseRoundingRadius)
    {
        translate([0,W/2-Rand,0]) cube([L/2-CaseRadius,Rand,H],center=false); // Wall
        translate([L/2-Rand,0,0]) cube([Rand,W/2-CaseRadius,H],center=false); // Wall
        translate([L/2-CaseRadius,W/2-CaseRadius,H/2]) cylinder(h=H,r=CaseRadius,center = true);
    }
    else
    {
        translate([0,W/2-Rand,0]) cube([L/2-Rad,Rand,H],center=false); // Wall
        translate([L/2-Rand,0,0]) cube([Rand,W/2-Rad,H],center=false); // Wall
    }
    translate(ScrewCornerPos) cylinder(h=H,r=Rad,center = false); // Cylinder
    translate([L/2-3*Rad+Rand,W/2-Rad,0]) rotate([0,0,  0]) HolderGap(H,Rad,Rand); // Gap between wall and Cylinder
    translate([L/2-Rad,W/2-Rad-Rand,0])   rotate([0,0,-90]) HolderGap(H,Rad,Rand); // Gap Between wall and Cylinder
    if (XAdditionalScrew)    {
        translate(ScrewAddXPos) cylinder(h=H,r=Rad,center = false); // Cylinder
        translate([Rand,W/2-Rad,0]) rotate([0,0,0]) HolderGap(H,Rad,Rand);
        translate([Rand-2*CaseRoundingRadius,W/2-Rad,0]) rotate([0,0,0]) HolderGap(H,Rad,Rand);
    }
    if (YAdditionalScrew)    {
        translate(ScrewAddYPos) cylinder(h=H,r=Rad,center = false); // Cylinder
        translate([L/2-3*Rad+2*CaseRoundingRadius,-Rand+2*CaseRoundingRadius,0]) rotate([0,0, 270]) HolderGap(H,Rad,Rand);
        translate([L/2-3*Rad+2*CaseRoundingRadius,-Rand,0]) rotate([0,0, 270]) HolderGap(H,Rad,Rand);
    }
}

module NutCut(TotalHigh,High,Dia){
    AdditionalGap=0.3;
    translate([0,0,-(High+2*AdditionalGap)/2-NutSink]) cylinder($fn=6,h=High+2*AdditionalGap,d=2*sqrt(((Dia/2)*(Dia/2))+((Dia/4)*(Dia/4)))+Dia/26+2*AdditionalGap,center = true);
    translate([0,0,-HoleDeepness/2]) cylinder(h=HoleDeepness,d=ScrewHoleDia,center = true);
}

module SquareNutCut (TotalHigh,High,Size,OnlyNut) {
    AdditionalGap=0.5;
    SquareNutInsertReduction= 0.2;
    if (OnlyNut)
    {
        translate([0,0,-(High+2*AdditionalGap)/2-NutSink])cube([Size+2*AdditionalGap,Size+2*AdditionalGap,High+2*AdditionalGap],center=true);
        translate([0,0,-HoleDeepness/2]) cylinder(h=HoleDeepness,d=ScrewHoleDia,center = true);
    }
    else
    {
        translate([0,0,-(High+2*AdditionalGap)/2-NutSink])cube([Size+2*AdditionalGap,Size+2*AdditionalGap,High+2*AdditionalGap],center=true);
        translate([CaseRoundingRadius/2+0.02,0,-(High+2*AdditionalGap)/2-NutSink+SquareNutInsertReduction/2]) cube([CaseRoundingRadius+0.04,Size+2*AdditionalGap,High+2*AdditionalGap-SquareNutInsertReduction],center=true);
        translate([0,0,-HoleDeepness/2]) cylinder(h=HoleDeepness,d=ScrewHoleDia,center = true);
    }
}

module SideWallHoles () {
    cylinder(h=20,d1=10,d2=15,center = true);
}

module EthernetPortCut (Width,Height,InsideDepth,Wall,Position,ZLocal) {
    // Width  = opening size along the wall, as seen from outside (e.g. 19mm)
    // Height = opening size vertically, as seen from outside (e.g. 17mm)
    // InsideDepth = how far the RJ45 module body extends past the inner wall face (e.g. 10mm)
    // ZLocal = vertical position in the calling body's own local coordinate frame (caller computes this)
    // Cut runs from the outer wall face through the wall thickness, plus InsideDepth further inward,
    // with a generous OL overlap margin on BOTH ends so it fully clears the collar (avoids a paper-thin
    // uncut membrane sealing the back of the hole)
    TotalDepth = SideWallThickness+InsideDepth;
    OL = 0.5;
    CutDepth = TotalDepth+2*OL;
    if (Wall=="A") { translate([Caselength/2-TotalDepth/2, Position, ZLocal]) cube([CutDepth,Width,Height],center=true); }
    if (Wall=="C") { translate([-Caselength/2+TotalDepth/2, Position, ZLocal]) cube([CutDepth,Width,Height],center=true); }
    if (Wall=="B") { translate([Position, -CaseWidth/2+TotalDepth/2, ZLocal]) cube([Width,CutDepth,Height],center=true); }
    if (Wall=="D") { translate([Position, CaseWidth/2-TotalDepth/2, ZLocal]) cube([Width,CutDepth,Height],center=true); }
}

module EthernetPortCollar (Width,Height,InsideDepth,CollarWall,Wall,Position,ZLocal) {
    // Solid sleeve of added material surrounding the port opening, spanning the InsideDepth zone
    // Modified to extend all the way down to the inner roof (Z = BottomTopThickness)
    // to prevent overhangs/printing in mid-air when the top lid is printed upside down.

    CW = Width+2*CollarWall;
    CH = Height+2*CollarWall;

    // Calculate a new height and Z-center that anchors the block to the lid plate
    NewHeight = (ZLocal + (CH/2)) - BottomTopThickness;
    NewZCenter = BottomTopThickness + (NewHeight / 2);

    if (Wall=="A") { translate([Caselength/2-SideWallThickness-InsideDepth/2, Position, NewZCenter]) cube([InsideDepth+0.02,CW,NewHeight],center=true); }
    if (Wall=="C") { translate([-Caselength/2+SideWallThickness+InsideDepth/2, Position, NewZCenter]) cube([InsideDepth+0.02,CW,NewHeight],center=true); }
    if (Wall=="B") { translate([Position, -CaseWidth/2+SideWallThickness+InsideDepth/2, NewZCenter]) cube([CW,InsideDepth+0.02,NewHeight],center=true); }
    if (Wall=="D") { translate([Position, CaseWidth/2-SideWallThickness-InsideDepth/2, NewZCenter]) cube([CW,InsideDepth+0.02,NewHeight],center=true); }
}

module RJ45Shape () {
    // Reference model of a physical RJ45 jack module, built along local +X = "depth", starting
    // at X=0 (the INNER wall face, i.e. the start of the friction-fit collar's InsideDepth zone -
    // NOT the outer face of the wall) and increasing inward. Width sits on local Y, centered on
    // Y=0.
    // Height sits on local Z: the BOTTOM face is held flat at Z = -NoseHeight/2 for the
    // entire model (nose through the taper tip); the TOP face is what slopes down in the taper
    // section. Z=0 (the nose's vertical center) still lines up with ZLocal, matching
    // EthernetPortCut/EthernetPortCollar's convention, since the nose itself stays symmetric.
    // Display-only - not unioned or subtracted from the printed bottom/top bodies, same pattern
    // as Breadboard()/BatteryHolder()/PCB1().
    //
    // Part 1 - nose: 19 x 17mm cross-section, exactly EthernetPortInsideDepth deep (the same
    // depth the collar itself occupies), so it sits perfectly inside the collar's InsideDepth
    // zone with nothing left protruding through the wall thickness itself.
    NoseWidth   = 19;
    NoseHeight  = 17;
    NoseDepth   = EthernetPortInsideDepth;

    // Part 2 - wider jack body, measured from the end of the nose (local X = NoseDepth):
    //   a) flare   : 0.0 -> 5.0mm   width 19 -> 34mm, smooth S-curve (raised-cosine / "sine
    //                                wave" - zero slope at both ends, steepest at the midpoint
    //                                "curve wall" at 2.5mm), height const 17mm, both faces
    //                                symmetric (Z centered).
    //   b) flat flange : 5.0 -> 7.0mm   width const 34mm, height const 17mm, Z still centered.
    //   c) height taper: 7.0 -> 11.0mm  width const 34mm. Bottom face stays flat the whole way;
    //                                    only the TOP face slopes down, taking the height from
    //                                    17mm to 12mm.
    FlareDepth  = 5;      // 0.25cm + 0.25cm, the two phases either side of the "curve wall"
    FlareSteps  = 24;     // slice count approximating the smooth curve (higher = smoother)
    FlangeWidth = 34;
    FlatDepth   = 2;
    TaperDepth  = 4;
    TaperHeight = 12;

    BottomZ    = -NoseHeight/2;       // fixed bottom face, held constant for the whole model
    TopZ_full  =  NoseHeight/2;       // top face while height is still the full 17mm
    TopZ_taper =  BottomZ + TaperHeight; // top face at the tip, after the one-sided taper

    // Thin slice spanning the FULL (untapered) height, at depth x, width w - used to loft the
    // width flare while height stays constant.
    module fullSlice(x,w) {
        translate([x,0,(BottomZ+TopZ_full)/2]) cube([0.01,w,(TopZ_full-BottomZ)],center=true);
    }
    // Thin slice at depth x, constant width FlangeWidth, with the bottom fixed at BottomZ and
    // top at topZ - used to loft the one-sided height taper.
    module taperSlice(x,topZ) {
        translate([x,0,(BottomZ+topZ)/2]) cube([0.01,FlangeWidth,(topZ-BottomZ)],center=true);
    }

    d0 = NoseDepth;                   // end of nose / start of flare
    d1 = d0 + FlareDepth;             // end of flare / start of flat flange
    d2 = d1 + FlatDepth;              // end of flat flange / start of height taper
    d3 = d2 + TaperDepth;             // end of height taper (model tip)

    // Front metal EMI shield plate, sized/placed per spec:
    //  - Depth (X): 0.1cm = 1.0mm, protruding forward past the nose's front face (X=0)
    //  - Height (Z): 1.3cm = 13.0mm, its TOP edge sits 0.1cm (1.0mm) below the nose's top face
    //    (TopZ_full) - the "height reduction" side, i.e. the side that later tapers down in Part 2c
    //  - Width (Y): 1.65cm = 16.5mm, centered on Y=0 (equal 1.25mm margin to each side of the 19mm nose)
    // NOTE: alpha is applied directly on each color() call here (Silver/DarkSlateGray) rather than
    // via a single enclosing color() wrapper around the whole shape in RJ45Model(). OpenSCAD's
    // OpenCSG *preview* (F5) does not correctly composite a child's own color() under an ancestor
    // color() - the ancestor's color silently wins in preview (a full F6/--render evaluates it
    // correctly). Since RJ45Model() previously wrapped this entire shape in one translucent dark
    // color, the plate was there geometrically, but preview rendered it identically to the rest of
    // the body - effectively invisible. Coloring each part here, with its own alpha, fixes it in
    // both preview and full render.
    PlateDepth   = 1.0;              // 0.1cm
    PlateHeight  = 13.0;             // 1.3cm
    PlateWidth   = 16.5;             // 1.65cm
    PlateGap     = 1.0;              // 0.1cm gap down from the nose's top face
    PlateOverlap = 0.01;             // tiny extra depth fused into the nose, avoids a coincident face


    // Standard RJ45 3-part plug profile based on custom specifications (converted to mm)
    // Part 1: Main body (1.3cm x 0.75cm)
    PortPart1_Width  = 13.0;
    PortPart1_Height = 7.5;

    // Part 2: Middle step (0.8cm x 0.2cm)
    PortPart2_Width  = 8.0;
    PortPart2_Height = 2.0;

    // Part 3: Top clip (0.55cm x 0.1cm)
    PortPart3_Width  = 5.5;
    PortPart3_Height = 1.0;

    // Calculate cut depths
    PortCutFrontX = -(PlateDepth + 0.5);
    PortCutBackX  = d0;
    PortCutDepth  = PortCutBackX - PortCutFrontX;
    PortCutCenterX = (PortCutFrontX + PortCutBackX)/2;

    // Calculate Z alignments to ensure the 0.2cm (2.0mm) distance from the bottom of the metal plate
    PlateCenterZ  = TopZ_full - PlateGap - PlateHeight/2;
    PlateBottomZ  = PlateCenterZ - PlateHeight/2;
    PortBaseZ     = PlateBottomZ + 2.0; // The 0.2cm gap at the non-reduced flat side

    difference() {
        union() {
            // Front Metal Plate
            color("Silver", 0.95)
            translate([-PlateDepth/2 + PlateOverlap/2, 0, TopZ_full - PlateGap - PlateHeight/2])
                cube([PlateDepth + PlateOverlap, PlateWidth, PlateHeight], center=true);

            // Group the rest of the plastic body in a distinct color so the metal plate remains visible
            color("DarkSlateGray", 0.85) {
                // Part 1: nose - straight extrusion
                translate([d0/2,0,(BottomZ+TopZ_full)/2]) cube([d0,NoseWidth,(TopZ_full-BottomZ)],center=true);

                // Part 2a: smooth S-curve flare (19mm -> 34mm)
                for (i = [0:FlareSteps-1]) {
                    t0 = i/FlareSteps;
                    t1 = (i+1)/FlareSteps;
                    x0 = d0 + FlareDepth*t0;
                    x1 = d0 + FlareDepth*t1;
                    w0 = NoseWidth + (FlangeWidth-NoseWidth)*(1-cos(180*t0))/2;
                    w1 = NoseWidth + (FlangeWidth-NoseWidth)*(1-cos(180*t1))/2;
                    hull() { fullSlice(x0,w0); fullSlice(x1,w1); }
                }

                // Part 2b: flat flange (34mm wide, 17mm tall, 2mm deep)
                translate([d1+FlatDepth/2,0,(BottomZ+TopZ_full)/2]) cube([FlatDepth,FlangeWidth,(TopZ_full-BottomZ)],center=true);

                // Part 2c: one-sided height taper
                hull() { taperSlice(d2,TopZ_full); taperSlice(d3,TopZ_taper); }
            }
        }

        // The standard 3-part port hollow itself
        union() {
            // Part 1 (Main body)
            // Added +0.02 to Z-height to create a 0.01mm upward overlap
            translate([PortCutCenterX, 0, PortBaseZ + PortPart1_Height/2])
                cube([PortCutDepth, PortPart1_Width, PortPart1_Height + 0.02], center=true);

            // Part 2 (Middle step)
            // Added +0.02 to Z-height to create 0.01mm downward and upward overlaps
            translate([PortCutCenterX, 0, PortBaseZ + PortPart1_Height + PortPart2_Height/2])
                cube([PortCutDepth, PortPart2_Width, PortPart2_Height + 0.02], center=true);

            // Part 3 (Top clip)
            // Added +0.02 to Z-height to create a 0.01mm downward overlap
            translate([PortCutCenterX, 0, PortBaseZ + PortPart1_Height + PortPart2_Height + PortPart3_Height/2])
                cube([PortCutDepth, PortPart3_Width, PortPart3_Height + 0.02], center=true);
        }
    }
}

module RJ45Model (Wall,Position,ZLocal) {
    // Places RJ45Shape() at a port location using the same Wall/Position/ZLocal convention as
    // EthernetPortCut()/EthernetPortCollar(), so it automatically lines up with whichever wall
    // and position the actual port cutout is configured for. The shape's local X=0 (start of the
    // nose) is placed at the INNER wall face - i.e. offset inward from the outer face by
    // SideWallThickness - which is exactly where EthernetPortCollar's InsideDepth zone begins.
    // This keeps the nose entirely inside the collar and clear of the wall material itself.
    if (ShowRJ45Model)
    {
        if (Wall=="C") { translate([-Caselength/2+SideWallThickness, Position, ZLocal]) RJ45Shape(); }
        if (Wall=="A") { translate([Caselength/2-SideWallThickness, Position, ZLocal]) mirror([1,0,0]) RJ45Shape(); }
        if (Wall=="B") { translate([Position, -CaseWidth/2+SideWallThickness, ZLocal]) rotate([0,0,90]) RJ45Shape(); }
        if (Wall=="D") { translate([Position, CaseWidth/2-SideWallThickness, ZLocal]) rotate([0,0,-90]) RJ45Shape(); }
    }
}

module ScrewCut(m,h,v){
// m = 3=M3  4=M4  5=M5 6=M6 usw...
// h = High of the screw inkl. head
// v = if screw head is to be sunk deeper

    ScrewHeadDia=m*2; // Berechnung des Schraubenkopf Durchmessers
    //ScrewCountersink=(m+8)/14-0.7; // Leichte ScrewCountersink damit Schraube nicht vorsteht
    ScrewHoleDia=m+1; // ScrewHoleDiadurchmesser

    translate([0,0,-0.01])  union(){ // Ganze Schraube

            translate([0,0,ScrewCountersink-0.01])cylinder( h = ScrewHeadDia/4,d1=ScrewHeadDia,d2=ScrewHeadDia/2,center=false); // Kegel (Abschraegung)
            translate([0,0,0]) cylinder( h = ScrewCountersink,d=ScrewHeadDia,center=false);  // ScrewCountersink
            // V7.11 PhuNguyenPT: guarded with v>0 - all 3 current call sites pass v=0 (no extra
            // sink depth requested), which previously still built a cylinder(h=0,...): a
            // zero-height solid. OpenSCAD's built-in cylinder() silently tolerated that as a
            // no-op; BOSL2's cylinder() (pulled in once roundedBox() uses BOSL2's cuboid())
            // correctly asserts h>0 and aborts the render. Skipping the call entirely when
            // v<=0 is the correct fix either way - a v=0 "Versenkung" is nothing to cut.
            if (v>0) translate([0,0,0.01])rotate([180,0,0])cylinder(h=v,d=ScrewHeadDia,center = false); // Versenkung
            translate([0,0,0.01])cylinder( h = h,d=ScrewHoleDia,center=false); //Loch fuer Gewinde
    }
}

module HolderGap (H,Rad,Rand) {
    difference(){
        translate([0,0,0]) cube([Rad*2-2*Rand,Rad-Rand,H],center=false);
        translate([0,0,-0.02]) cylinder(h=H+0.04,r=Rad-Rand,center = false);
        translate([2*(Rad-Rand),0,-0.02]) cylinder(h=H+0.04,r=Rad-Rand,center = false);
    }
}

module Breadboard () {
    // Semi-transparent reference block showing where the breadboard sits inside the case.
    // It is display-only - it is not unioned or subtracted from the printed bottom/top bodies.
    if (ShowBreadboard)
    {
        color([0.2,0.6,0.9,0.35])
        translate([BreadboardOffset_X,BreadboardOffset_Y,BottomTopThickness+BreadboardOffset_Z+BreadboardHeight/2])
            cube([BreadboardLength,BreadboardWidth,BreadboardHeight],center=true);
    }
}

module BatteryHolder () {
    // Illustrative reference model of the battery holder: hollow shell with a floor (mounting tab +
    // countersunk holes), outer perimeter walls, and 2 internal divider walls splitting the interior
    // into 3 equal battery slots. Display-only - not unioned or subtracted from the printed bottom/top bodies.
    if (ShowBatteryHolder)
    {
        // Battery holder mounting hole reference positions - same math used for the printed boss/thread
        // cuts in BodyBottom(), so these illustrative holes line up with the real bosses below.
        BatteryHolder_LeftX = BatteryHolderOffset_X - (BatteryHolderLength/2);
        BatteryHolder_RightX = BatteryHolderOffset_X + (BatteryHolderLength/2);
        BatteryHolder_BottomY = BatteryHolderOffset_Y - (BatteryHolderWidth/2);

        // The holder rests on TOP of the printed mounting bosses, not directly on the case floor
        BatteryHolderRestZ = BottomTopThickness + BatteryHolderBossHeight + BatteryHolderOffset_Z;

        // Inner cavity available for the 3 battery slots, measured between the two outer WIDTH walls
        BatteryHolderInnerWidth    = BatteryHolderWidth - 2*BatteryHolderWidthWallThickness;
        BatteryHolderSlotWidth     = (BatteryHolderInnerWidth - 2*BatteryHolderSlotWallThickness)/3;
        // Y-offset (from center) of each of the 2 internal divider walls
        BatteryHolderDividerOffset = BatteryHolderSlotWidth/2 + BatteryHolderSlotWallThickness/2;

        color([0.9,0.5,0.15,0.35])
        union() {
            // --- Floor (mounting tab) with 6 countersunk mounting holes ---
            // Local Z=0 (BatteryHolderRestZ) is the very bottom face of the holder, resting on the boss top.
            difference() {
                translate([BatteryHolderOffset_X,BatteryHolderOffset_Y,BatteryHolderRestZ+BatteryHolderFloorThickness/2])
                    cube([BatteryHolderLength,BatteryHolderWidth,BatteryHolderFloorThickness],center=true);

                for (i = [0 : 2]) {
                    translate([BatteryHolder_LeftX + BatteryHolderHoleGapFromWidthEdge, BatteryHolder_BottomY + BatteryHolderHoleGapFromLengthEdge + (i * BatteryHolderHoleSpacing), BatteryHolderRestZ]) {
                        translate([0,0,-0.01]) cylinder(h=BatteryHolderTabBoreDepth+0.02, d=BatteryHolderHoleDiameter, center=false);
                        translate([0,0,BatteryHolderTabBoreDepth-0.01]) cylinder(h=BatteryHolderTabCsinkDepth+0.02, r1=BatteryHolderHoleDiameter/2, r2=BatteryHolderCsinkTopRadius, center=false);
                    }
                    translate([BatteryHolder_RightX - BatteryHolderHoleGapFromWidthEdge, BatteryHolder_BottomY + BatteryHolderHoleGapFromLengthEdge + (i * BatteryHolderHoleSpacing), BatteryHolderRestZ]) {
                        translate([0,0,-0.01]) cylinder(h=BatteryHolderTabBoreDepth+0.02, d=BatteryHolderHoleDiameter, center=false);
                        translate([0,0,BatteryHolderTabBoreDepth-0.01]) cylinder(h=BatteryHolderTabCsinkDepth+0.02, r1=BatteryHolderHoleDiameter/2, r2=BatteryHolderCsinkTopRadius, center=false);
                    }
                }
            }

            // --- Outer WIDTH-direction side walls: thin in Y (0.2cm), run along X, taller (1.8cm) ---
            for (ySide = [-1, 1]) {
                translate([BatteryHolderOffset_X, BatteryHolderOffset_Y + ySide*(BatteryHolderWidth/2 - BatteryHolderWidthWallThickness/2), BatteryHolderRestZ+BatteryHolderFloorThickness+BatteryHolderWidthWallHeight/2])
                    cube([BatteryHolderLength, BatteryHolderWidthWallThickness, BatteryHolderWidthWallHeight], center=true);
            }

            // --- Outer LENGTH-direction end walls: thin in X (0.1cm), run along Y, shorter (1.1cm) ---
            for (xSide = [-1, 1]) {
                translate([BatteryHolderOffset_X + xSide*(BatteryHolderLength/2 - BatteryHolderLengthWallThickness/2), BatteryHolderOffset_Y, BatteryHolderRestZ+BatteryHolderFloorThickness+BatteryHolderLengthWallHeight/2])
                    cube([BatteryHolderLengthWallThickness, BatteryHolderWidth, BatteryHolderLengthWallHeight], center=true);
            }

            // --- 2 internal slot divider walls: split the interior into 3 equal battery slots ---
            // Same orientation and height as the width walls (run along X, parallel to the cell axis)
            for (k = [-1, 1]) {
                translate([BatteryHolderOffset_X, BatteryHolderOffset_Y + k*BatteryHolderDividerOffset, BatteryHolderRestZ+BatteryHolderFloorThickness+BatteryHolderWidthWallHeight/2])
                    cube([BatteryHolderLength, BatteryHolderSlotWallThickness, BatteryHolderWidthWallHeight], center=true);
            }
        }
    }
}

module PCB1 () {
    // Semi-transparent reference block showing the 45.0 x 40.0 x 10.1mm PCB resting on top of
    // Device Holder 1's 4 corner M2.5 standoffs. Display-only - it is not unioned or subtracted
    // from the printed bottom/top bodies; the real mechanical bosses/holes are the DeviceHolder()
    // calls in BodyBottom() above, which already use the same hole spacing (32.5 x 37.5mm).
    if (ShowPCB1)
    {
        // Top face of the printed standoff = where the PCB's underside actually rests
        PCB1_RestZ = BottomTopThickness + ScrewCylinderHeight1;

        color([0.1,0.45,0.25,0.55])
        difference() {
            // NOTE: PCB1_Width (40mm) goes on X, PCB1_Length (45mm) goes on Y - this matches
            // how the real standoffs are laid out (DeviceHolder_X_Distance1=32.5mm spans the
            // 40mm width edge, DeviceHolder_y_Distance1=37.5mm spans the 45mm length edge).
            // Swapping this would put the box 90 degrees off from the actual screw positions.
            translate([Offset_X_1, Offset_Y_1, PCB1_RestZ + PCB1_Height/2])
                cube([PCB1_Width, PCB1_Length, PCB1_Height], center=true);

            // 4 corner mounting holes - positioned exactly on the DeviceHolder1 screw axis
            // (DeviceHolder_X_Distance1 / DeviceHolder_y_Distance1) so they line up with the
            // real standoff holes underneath.
            for (xs = [-1,1]) {
                for (ys = [-1,1]) {
                    translate([DeviceHolder_X_Distance1/2*xs + Offset_X_1, DeviceHolder_y_Distance1/2*ys + Offset_Y_1, PCB1_RestZ - 0.01])
                        cylinder(h=PCB1_Height+0.02, d=PCB1_HoleDiameter, center=false);
                }
            }
        }
    }
}

module DeviceHolder (Col,CylHeight,CylDia,HoleDia,UseThread=false,ThreadSize=0,ThreadPitchVal=0,ThreadAngleVal=30,ThreadFitVal=0) {
    color(Col)translate([0,0,CylHeight/2+BottomTopThickness]) difference(){
        cylinder(h=CylHeight,d=CylDia,center = true);
        if (UseThread) {
            // Blind hole: 0.01 margin at the base (matches the corner-screw / battery-boss
            // convention, avoids a coincident face with the cylinder's own bottom cap).
            // Height is now CylHeight (not CylHeight-0.02), so the cut OVERSHOOTS the top
            // face by the same 0.01 the base is pulled back - this is what actually opens
            // the bore at the entry side. The previous CylHeight-0.02 version fell 0.01mm
            // SHORT of the top face on both ends, sealing the thread inside solid plastic
            // so no screw could ever reach it.
            translate([0,0,-CylHeight/2+0.01])
                ScrewThread(1.01*ThreadSize+1.25*ThreadFitVal, CylHeight, ThreadPitchVal, ThreadAngleVal, ThreadFitVal);
        } else {
            translate([0,0,0]) cylinder(h=CylHeight+0.05,d=HoleDia,center = true);
        }
    }
}

module pie(radius, angle, height, spin=0) {
    // Negative angles shift direction of rotation
    clockwise = (angle < 0) ? true : false;
    // Support angles < 0 and > 360
    normalized_angle = abs((angle % 360 != 0) ? angle % 360 : angle % 360 + 360);
    // Select rotation direction
    rotation = clockwise ? [0, 180 - normalized_angle] : [180, normalized_angle];
    // Render
    if (angle != 0) {
        rotate([0,0,spin]) linear_extrude(height=height)
        difference() {
            circle(radius);
            if (normalized_angle < 180) {
                union() for(a = rotation)
                    rotate(a) translate([-radius, 0, 0]) square(radius * 2);
            }
            else if (normalized_angle != 360) {
                intersection_for(a = rotation)
                    rotate(a) translate([-radius, 0, 0]) square(radius * 2);
            }
        }
    }
}

module roundedBox(size, radius, sidesonly) // Laenge, Breite, Hoehe, Radius, 0/1
{
    // ------------------------------------------------------------------------------------------
    // V7.11 PhuNguyenPT: re-implemented on top of BOSL2's cuboid() instead of the original
    // hand-rolled union of overlapping cubes + edge cylinders + corner spheres.
    //
    // WHY: the old version built a "rounded box" out of 3 overlapping cubes, up to 12 edge
    // cylinders and 8 corner spheres, unioned together. Each of those primitives is tangent to
    // its neighbours at an exact mathematical boundary (radius-size[axis]/2 etc.), so CGAL's
    // boolean engine frequently has to merge perfectly coincident/near-coincident faces at those
    // tangent seams. That is a classic source of "non-manifold edges" once the result is unioned
    // again with other shapes and then differenced with cutting cubes/cylinders - exactly what
    // MountHolder() does for every style (1-5), and especially styles 4/5 which union THREE
    // rounded boxes (straight + two 45 deg rotated copies) to build the arrow/diamond wall-mount
    // tab, then subtract corner cubes and a hulled oval hole from the result. With
    // EnableMountHolder = false none of this geometry exists, which is why the rest of the case
    // (which does not use roundedBox()) always exported clean.
    //
    // FIX: BOSL2's cuboid(rounding=...) builds the exact same *shape* (a box with rounded
    // vertical edges only, or a box with all 12 edges/8 corners rounded) but does so via a single
    // coherent, well-tested construction rather than a union of many tangent primitives, so it
    // does not leave degenerate/coincident faces behind. The call signature below (size, radius,
    // sidesonly) is unchanged, so MountHolder() and every other call site needs no edits.
    // Requires: include <BOSL2/std.scad> at the top of the file (see notes there), and the
    // ScrewCut() v>0 guard (see that module) so BOSL2's stricter cylinder() doesn't choke on the
    // file's pre-existing v=0 "no extra sink depth" calls.
    //
    //   sidesonly truthy  -> only the 4 vertical (Z-axis) edges rounded, flat top/bottom
    //                        (used by MountHolderStyle 2 and 3)
    //   sidesonly falsy   -> all 12 edges / 8 corners rounded (full 3D rounded box)
    //                        (used by MountHolderStyle 1, and internally by 4/5's diamond tab)
    // ------------------------------------------------------------------------------------------
    if (radius <= 0) {
        // No rounding requested - a plain cube is already manifold, skip cuboid() entirely.
        cube(size, center=true);
    }
    else if (sidesonly) {
        cuboid(size, rounding=radius, edges="Z", anchor=CENTER);
    }
    else {
        cuboid(size, rounding=radius, edges="ALL", anchor=CENTER);
    }
}

