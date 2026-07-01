// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Case
// Parametric OpenSCAD  |  v5.1 (Drop-In U-Slot & Pillars)
// Units: mm
// =====================================================

include <mke_s13_config.scad>

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
        translate([safe_line_x - 1, -5])
            square([box_l + 6, pcb_w + 10]);
    }
}

module shell_2d(w) {
    offset(r = w + pcb_gap) pcb_box_section_2d();
}

// =====================================================
// BOTTOM SHELL (Flat Top for Labyrinth)
// =====================================================
module bottom_shell() {
    difference() {
        // 1. Main Outer body box (Full height, flat top)
        linear_extrude(outer_h) shell_2d(wall);

        // 2. Interior hollow (Extruded slightly taller to ensure a clean boolean cut)
        translate([0, 0, floor_t])
            linear_extrude(outer_h + 1)
                offset(r = pcb_gap) pcb_box_section_2d();

        // --- DROP-IN U-SLOT FOR PCB ---
        translate([
            safe_line_x - 1 - pcb_gap - wall - 1.0,
            pcb_w/2 - (pcb_w/2 + slot_gap),
            floor_t + pcb_t + 0.2 - slot_gap
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
            offset(r = pcb_gap) pcb_box_section_2d();
    }
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
                offset(r = pcb_gap) pcb_box_section_2d();

        // --- DROP-IN U-SLOT FOR PCB (flared mouth) ---
        // Hull of a tight interior passage (unchanged clearance,
        // keeps the labyrinth seal reasonable) and a flared exterior
        // opening (extra mouth_chamfer_z in Z and Y) removes the
        // sharp sill edge that the cantilevered bare probe would
        // otherwise pivot against under insertion bending.
        hull() {
            // Exterior (soil-side) face -- flared
            translate([
                safe_line_x - 1 - pcb_gap - wall - 1.0,
                pcb_w/2 - (pcb_w/2 + slot_gap + mouth_chamfer_z),
                floor_t + pcb_t + 0.2 - slot_gap - mouth_chamfer_z
            ])
                cube([
                    0.01,
                    pcb_w + 2*(slot_gap + mouth_chamfer_z),
                    outer_h + 1 + mouth_chamfer_z
                ]);
            // Interior face -- tight, original clearance
            translate([
                safe_line_x - 1 - pcb_gap - wall - 1.0 + mouth_chamfer_x,
                pcb_w/2 - (pcb_w/2 + slot_gap),
                floor_t + pcb_t + 0.2 - slot_gap
            ])
                cube([
                    wall + 2.0 - mouth_chamfer_x,
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
            offset(r = pcb_gap) pcb_box_section_2d();
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
        translate([0, 0, -lip_h - 0.1])
            linear_extrude(lip_h + 0.2)
                shell_2d(0);

        // Cut the explicit labyrinth groove
        translate([0, 0, -lip_h - 0.01])
            linear_extrude(lip_h + 0.2)
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

    // --- CLOSURE BAFFLE (chamfered tip) ---
    baffle_h = inner_h - (2 * pcb_t + 0.2 + slot_gap + solder_z_tol);
    baffle_w = pcb_w + 2*slot_gap - baffle_clearance;

    hull() {
        translate([
            safe_line_x - 1 - pcb_gap - wall,
            pcb_w/2 - baffle_w/2,
            -(baffle_h - mouth_chamfer_z)
        ])
            cube([wall, baffle_w, baffle_h - mouth_chamfer_z + overlap_eps]);

        translate([
            safe_line_x - 1 - pcb_gap - wall + mouth_chamfer_x,
            pcb_w/2 - baffle_w/2 + mouth_chamfer_z,
            -baffle_h
        ])
            cube([wall - mouth_chamfer_x, baffle_w - 2*mouth_chamfer_z, 0.01]);
    }

    // --- CONNECTOR / PCB PARTITION BULKHEAD (UPPER HALF) ---
    partition_drop_h = outer_h - (z_pcb_seat + pcb_t + slot_gap);
    partition_w      = pcb_w + 2*pcb_gap - baffle_clearance;

    translate([
        partition_x1,
        pcb_w/2 - partition_w/2,
        -partition_drop_h
    ])
        cube([partition_t, partition_w, partition_drop_h + overlap_eps]);

    // --- INTEGRATED LOCKING PILLARS ---
    for (y_off = [-hole_sp/2, hole_sp/2]) {
        translate([hole_x, hole_cy + y_off, 0]) {
            shoulder_h = inner_h - (z_pcb_seat + pcb_t);
            translate([0, 0, -shoulder_h])
                cylinder(d = pcb_boss_d, h = shoulder_h + overlap_eps);

            pin_h = inner_h - shoulder_h - squish_tol;
            translate([0, 0, -inner_h + squish_tol])
                cylinder(d = lock_pin_d, h = pin_h);
        }
    }
}

// =====================================================
// PRINT-READY ORIENTATIONS
// =====================================================
// The lid() module above is authored in "assembly logic":
// the flat plate sits at the TOP (z=0..lid_t) and the
// closure baffle + locking pillars hang DOWNWARD into
// negative Z, matching how the lid sits when placed onto
// the bottom shell. That orientation is NOT printable as-is
// (parts fall below the z=0 bed plane), and printing it
// pillars-down would require support material under every
// pillar tip.
//
// lid_print_ready() mirrors the part so the plate becomes
// the base sitting flat on the bed, and the pillars become
// self-supporting upward-growing towers — no supports needed.
// Print the plate first; the pillars build up naturally on
// top of it, same logic as the bottom shell's standoffs.
module lid_print_ready() {
    translate([0, 0, lid_t])
        mirror([0, 0, 1])
            lid();
}

// =====================================================
// RENDER
// =====================================================
color("SteelBlue",  0.85) bottom_shell();
color("LightBlue",  0.70) translate([0, pcb_w + 15, 0]) lid_print_ready();