// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Shared Params
// Single source of truth for case + assembly files
// Units: mm
// =====================================================
// Included (not used) by mke_s13_case.scad and
// mke_s13_assembly.scad via: include <mke_s13_config.scad>
// =====================================================

// =====================================================
// 0. PRINT TOLERANCE PROFILE
// =====================================================
// This project is printed on a 0.2mm nozzle at 0.1mm layer height --
// finer than the 0.4mm nozzle most consumer FDM printers ship with
// by default (Ultimaker, Prusa, Bambu Lab, Creality, Anycubic, etc.),
// chosen here for the extra dimensional precision this design's
// small features (M3 threads, labyrinth lip, locking pins) benefit
// from. All radial clearances in this design (lock_pin_d,
// baffle_clearance) are sized to be comfortably
// larger than one nozzle width (>=0.75x nozzle_d/side), so ordinary
// FDM dimensional drift won't close the gap. If switching to a
// coarser/finer nozzle, revisit those clearance values relative to
// the actual nozzle_d in use -- and note the numbers embedded in
// their comments below (currently computed FOR nozzle_d=0.2) will
// need updating too, not just this value.
nozzle_d = 0.2;
clearance_per_side = nozzle_d * 0.75; // = 0.15mm at this project's
                                       // 0.2mm-nozzle setup -- the
                                       // per-side radial/lateral fit
                                       // tolerance used throughout this
                                       // design's sliding/drop-in joints
// =====================================================
// 0.5. REAL-WORLD MANUFACTURING TOLERANCES
// =====================================================
fab_x_tol    = 1.0; // 1. Extra X-axis length for PCB routing variations
solder_z_tol = 0.5; // 2. Extra Z-clearance under board for hand-solder spikes
squish_tol   = 0.2; // 3. Shortens locking pins to prevent elephant-foot bottoming out
sealant_tol  = 0.5; // 4. Extra gap around probe slot for silicone/conformal coating
overlap_eps  = 0.01; // Boolean-union safety margin: internal features that
                      // meet the lid plate at an exact z=0 seam (baffle,
                      // bulkhead, pillar shoulders) are extended this far
                      // past the seam so the union is a genuine 3D overlap
                      // rather than an exact coincident-face touch, which
                      // CGAL can occasionally resolve into a non-manifold
                      // edge. Purely a CAD-kernel safety margin -- has no
                      // effect on printed dimensions since it's buried
                      // inside material that's already solid there.
lid_center_hollow_overcut = 0.1;  // how far the lid's center hollow-out cut
                                   // reaches past z=0 into the plate (see lid()) --
                                   // must stay in sync with lid_reconnect_h below
lid_reconnect_h = lid_center_hollow_overcut + overlap_eps; // internal features
                                   // (baffle, bulkhead, pillar shoulders) that
                                   // attach to the underside of the lid plate must
                                   // extend at least this far past z=0, or they end
                                   // up as disconnected islands inside the void left
                                   // by the hollow-out cut
groove_overcut = 0.2; // Extra Z-depth for the labyrinth groove boolean
                      // cut
// =====================================================
// 1. PCB DIMENSIONS
// =====================================================
pcb_l       = 118.0; // total length (tip to right edge)
pcb_w       =  22.0; // width (parallel body section)
pcb_t       =   1.605; // thickness
chev_l      =   9.0; // length of chevron (pointed) section from tip
corner_r    =   2.0; // right-end corner radius

// =====================================================
// 2. SAFE LINE
// =====================================================
safe_line_x =  76.0; // Boundary between bare probe and case

// =====================================================
// 3. CONNECTOR — CASE CLEARANCE ENVELOPE (JST XH 3P, fully mated)
// =====================================================
conn_l      =  10.0; // length across pins (PCB width axis)
conn_h      =  10.0; // height above PCB top surface

// =====================================================
// 4. CONNECTOR — DETAILED FOOTPRINT (Assembly model, JST XH 2.5mm 3P Top-Entry)
// =====================================================
conn_pitch    = 2.54;   // pin pitch
conn_male_w   = 9.9;   // male shroud width
conn_male_d   = 5.75;  // male shroud depth
conn_male_h   = 7.0;   // male shroud height
pin_offset_x  = 2.35;  // pin offset from PCB edge

// =====================================================
// 5. PCB MOUNTING HOLES & INTERLOCKS
// =====================================================
hole_d      =   3.2; // PCB hole clearance diameter
hole_sp     =  15.0; // Y center-to-center spacing
hole_x      =  88.5; // X from tip
hole_cy     =  pcb_w / 2; // Y center (computed)
pcb_boss_d  =   6.7; // Outer diameter of standoffs/pillars
lock_pin_d  =   hole_d - 2*clearance_per_side; // Lid's locking-pillar pin
                      // diameter (was hardcoded inline as 3.0mm -- 0.1mm/side
                      // clearance vs hole_d was too tight for FDM; now derived
                      // from clearance_per_side, giving 0.15mm/side at this
                      // project's actual 0.2mm-nozzle setup)

// =====================================================
// 6. CASE WALLS & CLEARANCES
// =====================================================
wall        =   2.0; // wall thickness
floor_t     =   2.0; // floor thickness
lid_t       =   2.0; // lid plate thickness
pcb_gap_y     =   clearance_per_side; // PCB-to-inner-wall clearance
shoulder_pcb_clearance_z = 0.4;
pcb_z_gap = 0.2;
slot_gap    =   clearance_per_side + sealant_tol; // extra clearance around probe
                                                  // passthrough
baffle_clearance = 2 * clearance_per_side + sealant_tol; // diametral width clearance
                         // between closure baffle and the U-slot it drops
                         // into (was a hardcoded -0.4 giving 0.2mm/side; now
                         // derived from clearance_per_side, giving 0.15mm/side
                         // at this project's actual 0.2mm-nozzle setup)

// =====================================================
// 7. HEIGHT STACK (all computed)
// =====================================================
// stack_h is measured top-down: from the top of the mated
// male+female connector, through the PCB, to the tip of the
// solder pins protruding below the PCB's bottom surface.
// Confirmed by physical measurement: mated connector height
// (conn_h) = 10.0mm, total stack = 12.115mm.
stack_h         =  12.115;            // connector top -> pin tip ends (measured)
pin_protrusion  =  stack_h - conn_h - pcb_t; // solder pin tip protrusion
                                             // below the PCB bottom surface
                                             // (derived so stack_h stays
                                             // internally consistent)
inner_h     =  stack_h + solder_z_tol; // interior clear height
outer_h     =  inner_h + floor_t;  // total shell height
z_pcb_seat  =  floor_t + pcb_t + pcb_z_gap + solder_z_tol; // PCB seat height inside the case

// =====================================================
// 8. CABLE EXIT CLEARANCE
// =====================================================
cable_clear =   1.0; // clearance around connector footprint

// =====================================================
// 9. DERIVED LAYOUT
// =====================================================
box_l       =  pcb_l - safe_line_x; // enclosure length (currently 42mm)

// =====================================================
// 10. CONNECTOR / PCB PARTITION BULKHEAD
// =====================================================
// All four values below are measured from the PCB's RIGHT edge
// (x = pcb_l), i.e. the end opposite the chevron tip, running
// along the PCB length axis -- matching how this was physically
// measured on the board:
//   0.0 -- 0.6 cm : connector footprint (JST shrouds sit here,
//                   confirmed against connector_male() in the
//                   assembly file: shroud starts ~112.25mm)
//   0.7 -- 0.85cm : bulkhead wall (this section) -- separates the
//                   connector cavity from the main PCB cavity
//   0.9 -- 2.6 cm : PCB "red line" keep-out zone (informational;
//                   no case feature is placed here, but it's why
//                   the bulkhead is pinned at 0.85cm and not pushed
//                   further left)
// NOTE: the bulkhead is built into lid(), NOT bottom_shell(). The
// PCB and male connector are pre-soldered into one rigid unit
// before assembly, so the bottom shell must stay a plain,
// unobstructed cavity for that unit to drop straight into. The
// bulkhead only comes down afterward with the lid, around the
// already-seated PCB.
partition_far_edge  =   7.0;  // mm from right edge -- wall face nearer the connectors
partition_near_edge =   8.5;  // mm from right edge -- wall face nearer the tip/PCB body
red_line_near_edge  =   9.0;  // mm from right edge -- start of PCB red-line zone (reference only)
red_line_far_edge   =  26.0;  // mm from right edge -- end of PCB red-line zone (reference only)

partition_x1 = pcb_l - partition_near_edge; // = 109.5 -- tip-side face
partition_x2 = pcb_l - partition_far_edge;  // = 111.0 -- connector-side face
partition_t  = partition_x2 - partition_x1; // = 1.5mm -- bulkhead wall thickness
                                             // (independent of lid_t; derived from
                                             // physical PCB connector/keepout layout)

// =====================================================
// LABYRINTH JOINT PARAMETERS
// =====================================================
lip_h         = 3.0;     // Height of the tongue/groove overlap
lip_clear     = 0.15;    // Radial clearance for a smooth sliding fit

// Explicit, sturdy thicknesses for the lid walls:
inner_skirt_t = 1.0;     // Gives the lid's inner lip 1mm of solid plastic
outer_skirt_t = 2.0;     // Gives the lid's outer lip 2.0mm of solid plastic (overhang)

// The tongue on the base is automatically calculated to sit perfectly between them:
tongue_in  = inner_skirt_t + lip_clear;
tongue_out = wall;       // Tongue extends flush to the outer edge of the bottom shell

// =====================================================
// PROBE MOUTH CHAMFER (Cantilever stress-riser mitigation)
// =====================================================
// The bare probe (chevron tip to safe_line_x, ~76mm of unsupported
// 1.6mm FR4) exits the case through the U-slot / closure baffle gap.
// Under insertion-force bending, the PCB pivots against whichever
// edge it contacts first -- the top of the bottom_shell sill, or the
// bottom-front tip of the lid's closure baffle. Left as plain cube
// corners, those are hard 90 deg stress risers sitting right at the
// point of maximum bending moment -- the most likely crack-initiation
// site on the whole assembly. Flaring both into a funnel spreads
// contact over a slope instead of a line/edge.
 boolean_overlap = 1.0;