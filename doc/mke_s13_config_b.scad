// =====================================================
// MKE-S13 Variant B — Shared Overrides
// Single source of truth for everything mke_s13_case_b.scad
// changes from the base config. Both mke_s13_case_b.scad and
// mke_s13_assembly_b.scad `include` (not `use`) this file, so
// they always see identical values -- no manual re-declaring,
// no drift.
// =====================================================
include <mke_s13_config.scad>

// ---- Variant B PCB gap overrides (case expansion fix) ----
pcb_gap_x_left  = 2.5;  // exit-wall (probe side) -- sized for the exit bosses
pcb_gap_x_right = 4.5;  // connector-wall (opposite the pointy tip)
pcb_gap_y_left  = 12;   // exit-side long-edge gap -- sized for exit bosses
pcb_gap_y_right = 12;   // connector-side long-edge gap -- sized for right-corner bosses
lid_t = 3.5;

// ---- Screw boss / clearance geometry (variant B only) ----
// ---- Screw boss / clearance geometry (variant B only) ----
screw_clear_d            = 3.4;   // M3 clearance hole dia through the lid
screw_head_d             = 6.0;   // flat-head (countersunk) max dia
screw_countersink_depth  = 1.7;   // depth of the conical countersink
base_plug                = 2.0;   // solid plastic left under the thread cavity
// Added to the M3 nominal diameter (3.0mm) before cutting the internal
// thread with metric_thread(internal=true). 0.0 cuts the thread at exact
// ISO nominal, which FDM printers reliably under-size.
//
// Derived from clearance_per_side (Section 0) rather than a standalone
// number -- that section already exists as this project's single source
// of truth for nozzle-driven print tolerance, and its header comment
// explicitly names M3 threads as one of the features it's meant to serve.
// Hardcoding a separate value here would just be the same "duplicate
// tolerance logic" mistake as the old case_b.scad parameter redeclaration.
//
// Doubled (not used bare) to match how clearance_per_side is applied
// elsewhere for DIAMETRAL fits (lock_pin_d, baffle_clearance): a printed
// bore shrinks inward around its full circumference (crest rounding,
// extrusion bulge, shrinkage), the same all-sides physical picture as a
// pin-in-hole fit, not a one-sided wall gap like pcb_gap_y/slot_gap. At
// this project's nozzle_d=0.2 that's 2 * 0.15 = 0.30mm.
//
// Threaded bores aren't a perfect analogue to a plain clearance gap --
// crest/root geometry has its own printability limits beyond straight
// radial shrinkage -- so treat this as a principled starting point, not
// a proven number. PRINT A SINGLE TEST BOSS + SCREW before committing to
// a full print. If the coupon comes out too loose, dial back toward
// 1 * clearance_per_side rather than picking a new unrelated constant.
thread_fit_comp          = 2 * clearance_per_side;

boss_h      = outer_h;            // absolute Z, boss stands on the floor
pilot_depth = boss_h - base_plug;

// ---- 4-corner screw boss/hole positions (single source of truth) ----
// See mke_s13_case_b.scad's original comments for the full rationale on
// why left and right use different (linear vs. diagonal) constructions --
// unchanged here, just relocated so both files see the exact same result.
base_x_left           = safe_line_x - boolean_overlap;
interior_wall_x_left  = base_x_left - pcb_gap_x_left;
x_off_left            = boss_radius + boss_wall_gap_x;
x_left                = interior_wall_x_left + x_off_left;
y_off_left            = boss_radius + boss_wall_gap_y;
y_bot_L = 0 - y_off_left;
y_top_L = pcb_w + y_off_left;

inner_r_right_x = pcb_gap_x_right + corner_r;
inner_r_right_y = pcb_gap_y_right + corner_r;
dist_right_x = inner_r_right_x - boss_radius - boss_wall_gap_x;
dist_right_y = inner_r_right_y - boss_radius - boss_wall_gap_y;
offset_R_x = dist_right_x / sqrt(2);
offset_R_y = dist_right_y / sqrt(2);

base_x_right = (pcb_l + fab_x_tol) - corner_r;
x_right      = base_x_right + offset_R_x;
y_bot_R      = corner_r - offset_R_y;
y_top_R      = (pcb_w - corner_r) + offset_R_y;

corner_positions = [
    [x_left,  y_bot_L], // Left-Bottom
    [x_left,  y_top_L], // Left-Top
    [x_right, y_bot_R], // Right-Bottom
    [x_right, y_top_R]  // Right-Top
];