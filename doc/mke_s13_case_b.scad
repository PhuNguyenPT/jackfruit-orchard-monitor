// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Case (Variant B)
// 4-corner M3 screw bosses for lid retention
// Parametric OpenSCAD  | v2.3 (Corrected Screw Boss Positioning)
// Units: mm
// =====================================================

include <mke_s13_config_b.scad>
use <threads.scad>

// --- VIEW CONFIGURATION ---
// 1 = Print layout   (bottom shell + print-ready lid, side by side, bed-flat)
// 2 = Assembled      (bottom shell + lid closed, with transparency, CASE ONLY)
// 3 = Cutaway        (sliced along Y-midplane -- also shows the corner
//                      screw bosses/holes in section, CASE ONLY)
// 4 = Exploded       (bottom shell + lid separated vertically, CASE ONLY)
view_mode = 3;

$fn = 48;

// Left Side: NOT a real PCB corner. The PCB edge here is a plain straight
// line (y=0 / y=pcb_w run straight from the chevron tip all the way past
// this point) -- the "corner" at (base_x_left, 0)/(base_x_left, pcb_w) is
// an ARTIFICIAL 90-degree corner created only by intersecting the box
// section with a square boundary at x = safe_line_x - boolean_overlap, for
// the purposes of building the shell profile. It does not describe the
// true PCB shape, so a diagonal/tangent-to-corner formula there quietly
// under-clears the boss from the PCB: at pcb_gap_y=9, boss_d=7, gap=0.3, the
// old sqrt(2) construction left only ~0.18mm of real clearance to the PCB
// edge (not the intended 0.3mm) -- close enough to bind against ordinary
// FDM dimensional variance. Fix: treat X and Y as fully independent LINEAR
// offsets measured directly off their own real reference edges (the case's
// front boundary at x = safe_line_x - boolean_overlap for X, and the PCB's
// actual straight long edges at y=0 / y=pcb_w for Y). No sqrt(2) term --
// that only belongs to a genuine rounded corner, which isn't present here.
base_x_left           = safe_line_x - boolean_overlap;
interior_wall_x_left  = base_x_left - pcb_gap_x_left;
x_off_left             = boss_radius + boss_wall_gap_x;
x_left                 = interior_wall_x_left + x_off_left;
y_off_left              = boss_radius + boss_wall_gap_y;
y_bot_L = 0 - y_off_left;
y_top_L = pcb_w + y_off_left;

// Right Side: Base shape already has 'corner_r' (2.0mm).
// offset() expands this corner into an arc of radius (pcb_gap_y + corner_r).
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
module screw_boss(x, y) {
    translate([x, y, 0]) {
        difference() {
            // The solid outer boss remains the same
            cylinder(d = boss_d, h = boss_h);

            // The threaded internal cutout
            translate([0, 0, boss_h - pilot_depth])
                render()
                metric_thread(
                    diameter=3.0 + thread_fit_comp,
                    pitch=0.5,
                    length=pilot_depth + 1,
                    internal=true
                );
        }
    }
}

module screw_bosses() {
    for (p = corner_positions)
        screw_boss(p[0], p[1]);
}

// Lid clearance hole + shallow top chamfer at one corner.
boss_top_local = boss_h - outer_h; // negative -- boss top, in lid-local Z

module screw_clearance_hole(x, y) {
    translate([x, y, boss_top_local - 0.5])
        cylinder(d = screw_clear_d, h = (lid_t - boss_top_local) + 1);
}

module screw_countersink(x, y) {
    translate([x, y, lid_t - screw_countersink_depth])
        cylinder(d1 = screw_clear_d, d2 = screw_head_d, h = screw_countersink_depth + 0.05);
}

module screw_clearance_features() {
    for (p = corner_positions) {
        screw_clearance_hole(p[0], p[1]);
        screw_countersink(p[0], p[1]);
    }
}

// =====================================================
// 2D PROFILES
// =====================================================
module pcb_2d() {
    hull() {
        translate([0,        pcb_w/2])         circle(r=0.1);
        translate([chev_l,   0])                circle(r=0.1);
        translate([chev_l,   pcb_w])           circle(r=0.1);
        translate([pcb_l - corner_r + fab_x_tol, corner_r])           circle(r=corner_r);
        translate([pcb_l - corner_r + fab_x_tol, pcb_w - corner_r])   circle(r=corner_r);
    }
}

module pcb_box_section_2d() {
    intersection() {
        pcb_2d();
        translate([safe_line_x - boolean_overlap, -5])
            square([box_l + 6, pcb_w + 10]);
    }
}

// Directional "offset": scale-offset-unscale trick to expand rx along X, ry along Y
// instead of OpenSCAD's radius-only offset().
module aniso_offset(rx, ry) {
    scale([1, ry / rx, 1])
        offset(r = rx)
        scale([1, rx / ry, 1])
        children();
}

module shell_2d(w) {
    split_x = partition_x1;
    union() {
        intersection() {
            aniso_offset(w + pcb_gap_x_left, w + pcb_gap_y_left) pcb_box_section_2d();
            translate([-1000, -1000]) square([1000 + split_x, 2000]);
        }
        intersection() {
            aniso_offset(w + pcb_gap_x_right, w + pcb_gap_y_right) pcb_box_section_2d();
            translate([split_x, -1000]) square([1000, 2000]);
        }
    }
}

// =====================================================
// BOTTOM SHELL (Upstanding Tongue for Labyrinth - NO BEAD)
// =====================================================
module bottom_shell() {
    difference() {
        union() {
            linear_extrude(outer_h - lip_h) shell_2d(wall);
            translate([0, 0, outer_h - lip_h])
                linear_extrude(lip_h)
                    difference() {
                        shell_2d(tongue_out);
                        shell_2d(tongue_in);
                    }
        }

        translate([0, 0, floor_t])
            linear_extrude(outer_h + 1)
                shell_2d(0);

        hull() {
            translate([
                safe_line_x - boolean_overlap - pcb_gap_x_left - wall - 1.0,
                pcb_w/2 - (pcb_w/2 + slot_gap),
                z_pcb_seat - slot_gap
            ])
                cube([
                    0.01,
                    pcb_w + 2*(slot_gap),
                    outer_h + 1
                ]);
            translate([
                safe_line_x - boolean_overlap - pcb_gap_x_left - wall - 1.0,
                pcb_w/2 - (pcb_w/2 + slot_gap),
                z_pcb_seat - slot_gap // flat -- matches outer face above
            ])
                cube([
                    wall + 2.0,
                    pcb_w + 2*slot_gap,
                    outer_h + 1
                ]);
        }
    }

    // PCB mounting standoffs
    standoff_h = z_pcb_seat - floor_t;
    translate([0, 0, floor_t]) {
        for (y_off = [-hole_sp/2, hole_sp/2]) {
            translate([hole_x, hole_cy + y_off, 0])
            difference() {
                cylinder(d = pcb_boss_d, h = standoff_h);
                translate([0, 0, -0.01])
                    cylinder(d = hole_d, h = standoff_h + 1);
            }
        }
    }

    // --- CONNECTOR / PCB PARTITION BULKHEAD (LOWER HALF) ---
    intersection() {
        translate([partition_x1, -1000, 0])
            cube([partition_t, 2000, z_pcb_seat]);
        linear_extrude(z_pcb_seat)
            shell_2d(0);
    }

    // --- 4-CORNER SCREW BOSSES (variant B) ---
    screw_bosses();
}

// =====================================================
// LID (Double-Wall Groove Receiver - NO BEAD)
// =====================================================
module lid() {
    difference() {
        union() {
            linear_extrude(lid_t)
                shell_2d(wall + lip_clear + outer_skirt_t);

            translate([0, 0, -lip_h])
                linear_extrude(lip_h)
                    shell_2d(wall + lip_clear + outer_skirt_t);
        }

        translate([0, 0, -lip_h - lid_center_hollow_overcut])
            linear_extrude(lip_h + 2*lid_center_hollow_overcut)
                shell_2d(0);

        translate([0, 0, -lip_h - 0.01])
            linear_extrude(lip_h + groove_overcut)
                difference() {
                    shell_2d(tongue_out + lip_clear);
                    shell_2d(tongue_in - lip_clear);
                }

        translate([
            partition_x2,
            pcb_w/2 - (conn_l / 2) - cable_clear,
            -lip_h - 0.1
         ])
            cube([
                (pcb_l + cable_clear) - partition_x2,
                conn_l + 2*cable_clear,
                lid_t + lip_h + 0.2
            ]);

        // --- 4-CORNER SCREW CLEARANCE HOLES (variant B) ---
        screw_clearance_features();
    }

    baffle_h = inner_h - (2 * pcb_t + 0.2 + slot_gap + solder_z_tol);
    baffle_w = pcb_w + 2*slot_gap - (2 * clearance_per_side);

    hull() {
        // 1. Outer Face (Wide, fills the flared exterior cutout in X/Y only)
        translate([
            safe_line_x - boolean_overlap - pcb_gap_x_left - wall,
            pcb_w/2 - baffle_w/2 + clearance_per_side,
            -baffle_h
        ])
            cube([
                0.01,
                baffle_w - 2*clearance_per_side,
                baffle_h + lid_reconnect_h + groove_overcut
            ]);

        // 2. Inner Body (Narrow, matches interior cutout in X/Y only)
        translate([
            safe_line_x - boolean_overlap - pcb_gap_x_left - wall,
            pcb_w/2 - baffle_w/2,
            -baffle_h // flat -- matches outer face above
        ])
            cube([
                wall,
                baffle_w,
                baffle_h + lid_reconnect_h + groove_overcut
            ]);
    }

    partition_drop_h = outer_h - (z_pcb_seat + pcb_t + slot_gap);
    partition_w      = pcb_w + 2*pcb_gap_y_right - baffle_clearance;

    translate([
        partition_x1,
        pcb_w/2 - partition_w/2,
        -partition_drop_h
    ])
        cube([partition_t, partition_w, partition_drop_h + lid_reconnect_h]);

    for (y_off = [-hole_sp/2, hole_sp/2]) {
        translate([hole_x, hole_cy + y_off, 0]) {
             shoulder_h = outer_h - (z_pcb_seat + pcb_t) - shoulder_pcb_clearance_z;
            translate([0, 0, -shoulder_h])
                cylinder(d = pcb_boss_d, h = shoulder_h + lid_reconnect_h);

            pin_h = inner_h - shoulder_h - squish_tol;
            translate([0, 0, -inner_h + squish_tol])
                cylinder(d = lock_pin_d, h = pin_h);
        }
    }
}

module lid_print_ready() {
    translate([0, 0, lid_t])
        mirror([0, 0, 1])
            lid();
}

// =====================================================
// VIEW-MODE RENDERING MODULES (CASE ONLY -- no PCB/connectors)
// =====================================================
module print_layout() {
    plate_gap = 2.0;
    y_gap_max = max(pcb_gap_y_left, pcb_gap_y_right);
    bottom_shell_max_y = pcb_w + 2*(wall + y_gap_max);
    lid_min_y_offset    = wall + lip_clear + outer_skirt_t + y_gap_max;

    color("SteelBlue",  0.85) bottom_shell();
    color("LightBlue",  0.70)
        translate([0, bottom_shell_max_y + lid_min_y_offset + plate_gap, 0])
            lid_print_ready();
}

module full_system() {
    color("SteelBlue", 0.65) bottom_shell();
    translate([0, 0, outer_h]) color("LightBlue", 0.50) lid();
}

module exploded_system() {
    gap = 25;

    color("SteelBlue", 0.85) bottom_shell();
    translate([0, 0, outer_h + gap]) color("LightBlue", 0.75) lid();
}

// =====================================================
// RENDER EXECUTION
// =====================================================
if (view_mode == 1) {
    print_layout();
} else if (view_mode == 2) {
    full_system();
} else if (view_mode == 3) {
    difference() {
        full_system();
        translate([-10, pcb_w/2, -5]) {
            cube([pcb_l + 30, pcb_w + 20, outer_h + 30]);
        }
    }
} else if (view_mode == 4) {
    exploded_system();
}