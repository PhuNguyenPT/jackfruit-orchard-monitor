// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Case
// Parametric OpenSCAD  |  v6.1 (Case-Only View Modes)
// Units: mm
// =====================================================

include <mke_s13_config.scad>

// --- VIEW CONFIGURATION ---
// 1 = Print layout   (bottom shell + print-ready lid, side by side, bed-flat)
// 2 = Assembled      (bottom shell + lid closed, with transparency, CASE ONLY)
// 3 = Cutaway        (sliced along Y-midplane to verify clearances, CASE ONLY)
// 4 = Exploded       (bottom shell + lid separated vertically, CASE ONLY)
view_mode = 1;

$fn = 48;

// =====================================================
// 2D PROFILES
// =====================================================
module pcb_2d() {
    hull() {
        translate([0,        pcb_w/2])         circle(r=0.1);
        translate([chev_l,   0])                circle(r=0.1);
        translate([chev_l,   pcb_w])           circle(r=0.1);
        // X-axis extended to ensure rapid-fab PCBs don't crash into the wall
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

module shell_2d(w) {
    offset(r = w + pcb_gap_y) pcb_box_section_2d();
}

// =====================================================
// BOTTOM SHELL (Upstanding Tongue for Labyrinth - NO BEAD)
// =====================================================
module bottom_shell() {
    difference() {
        union() {
            // 1. Main Outer body box (Stops early to create the outer shoulder)
            linear_extrude(outer_h - lip_h) shell_2d(wall);

            // 2. Upstanding Centered Tongue (Clean, straight wall)
            translate([0, 0, outer_h - lip_h])
                linear_extrude(lip_h)
                    difference() {
                        shell_2d(tongue_out);
                        shell_2d(tongue_in);
                    }
        }

        // 3. Interior hollow (Leaves the internal cavity perfectly clear)
        translate([0, 0, floor_t])
            linear_extrude(outer_h + 1)
                offset(r = pcb_gap_y) pcb_box_section_2d();

        // --- DROP-IN SLOT FOR PCB (flat, rectangular -- slot_gap clearance
        // on all 4 sides, referenced directly off z_pcb_seat so it's exactly
        // slot_gap below the PCB's actual seated bottom surface, not a
        // reconstruction that silently drops solder_z_tol) ---
        translate([
            safe_line_x - boolean_overlap - pcb_gap_y - wall - 1.0,
            -slot_gap,
            z_pcb_seat - slot_gap
        ])
            cube([
                wall + 2.0,
                pcb_w + 2*slot_gap,
                outer_h + 1
            ]);
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
            offset(r = pcb_gap_y) pcb_box_section_2d();
    }
}

// =====================================================
// LID (Double-Wall Groove Receiver - NO BEAD)
// =====================================================
module lid() {
    // 1. CREATE MAIN SHELL & HOLLOW IT OUT
    difference() {
        union() {
            // Main Lid plate
            linear_extrude(lid_t)
                shell_2d(wall + lip_clear + outer_skirt_t);
            // Solid outer rim extending downwards
            translate([0, 0, -lip_h])
                linear_extrude(lip_h)
                    shell_2d(wall + lip_clear + outer_skirt_t);
        }

        // Hollow out the very center
        translate([0, 0, -lip_h - lid_center_hollow_overcut])
            linear_extrude(lip_h + 2*lid_center_hollow_overcut)
                shell_2d(0);

        // Cut the explicit labyrinth groove
        translate([0, 0, -lip_h - 0.01])
            linear_extrude(lip_h + groove_overcut)
                difference() {
                    shell_2d(tongue_out + lip_clear);
                    shell_2d(tongue_in - lip_clear);
                }

        // TIGHT CONNECTOR SLOT
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
    } // <-- THE DIFFERENCE BLOCK ENDS HERE

    // 2. ADD INTERNAL FEATURES (Now safe from being hollowed out)

    // --- CLOSURE BAFFLE (flat, rectangular) ---
    // baffle_h is the exact drop from the lid mating plane (outer_h) down to
    // slot_gap above the PCB's actual top surface (z_pcb_seat + pcb_t):
    //   outer_h - (z_pcb_seat + pcb_t + slot_gap)
    //     = (inner_h + floor_t) - (floor_t + pcb_t + pcb_z_gap + solder_z_tol + pcb_t + slot_gap)
    //     = inner_h - (2*pcb_t + pcb_z_gap + slot_gap + solder_z_tol)
    baffle_h = inner_h - (2 * pcb_t + pcb_z_gap + slot_gap + solder_z_tol);
    // Narrower than the shell's slot by clearance_per_side/side, purely so the
    // baffle drops into the shell's slot without binding on its side walls --
    // unrelated to PCB clearance, which is handled by baffle_h alone.
    baffle_w = pcb_w + 2*slot_gap - (2 * clearance_per_side);

    translate([
        safe_line_x - boolean_overlap - pcb_gap_y - wall,
        pcb_w/2 - baffle_w/2,
        -baffle_h
    ])
        cube([
            wall,
            baffle_w,
            baffle_h + lid_reconnect_h + groove_overcut
        ]);

    // --- CONNECTOR / PCB PARTITION BULKHEAD (UPPER HALF) ---
    partition_drop_h = outer_h - (z_pcb_seat + pcb_t + slot_gap);
    partition_w      = pcb_w + 2*pcb_gap_y - baffle_clearance;

    translate([
        partition_x1,
        pcb_w/2 - partition_w/2,
        -partition_drop_h
    ])
        cube([partition_t, partition_w, partition_drop_h + lid_reconnect_h]);

    // --- INTEGRATED LOCKING PILLARS ---
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

// =====================================================
// PRINT-READY ORIENTATIONS
// =====================================================
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
    bottom_shell_max_y = pcb_w + 2*(wall + pcb_gap_y);
    lid_min_y_offset    = wall + lip_clear + outer_skirt_t + pcb_gap_y;

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