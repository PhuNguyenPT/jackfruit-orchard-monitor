// =====================================================
// MKE-S13 Case + PCB — Combined Fit-Check Assembly
// Parametric OpenSCAD  |  v8.0 (Corrected PCB Triangle Tip)
// Units: mm
// =====================================================

use <mke_s13_case.scad>
include <mke_s13_config.scad>

// --- VIEW CONFIGURATION ---
// 1 = Assembled (Closed with transparency)
// 2 = Exploded (Separated vertically)
// 3 = Cutaway (Sliced along Y-midplane to verify clearances)
view_mode = 3;

$fn = 64;

// =====================================================
// PCB SUB-COMPONENTS
// =====================================================
module pcb_outline_2d() {
    hull() {
        // Point of the triangle
        translate([0, pcb_w/2]) circle(r=0.1);

        // Base of the triangle (where the rectangle begins, 9.0mm from the point)
        translate([chev_l, 0]) circle(r=0.1);
        translate([chev_l, pcb_w]) circle(r=0.1);

        // Far end of the rectangle
        translate([pcb_l - corner_r, corner_r]) circle(r=corner_r);
        translate([pcb_l - corner_r, pcb_w - corner_r]) circle(r=corner_r);
    }
}

module pcb_board() {
    color("ForestGreen", 0.95) {
        difference() {
            // Extrude the updated two-part hull (Rectangle + 0.9cm Triangle)
            linear_extrude(pcb_t) pcb_outline_2d();

            // 3 connector through-holes
            for(i = [-1, 0, 1]) {
                translate([pcb_l - pin_offset_x, pcb_w/2 + (i * conn_pitch), -0.5]) {
                    cylinder(d=1.0, h=pcb_t + 1);
                }
            }

            // 2 mounting/locking-pillar through-holes
            // (lid's interlocking pin passes through here into the
            // bottom standoff bore — see lid() in mke_s13_case.scad)
            for (y_off = [-hole_sp/2, hole_sp/2]) {
                translate([hole_x, hole_cy + y_off, -0.5]) {
                    cylinder(d=hole_d, h=pcb_t + 1);
                }
            }
        }
    }
}

module connector_male(z_extra = 0) {
    // 1. Male Header Plastic Shroud
    color("White", 0.95) {
        translate([pcb_l - conn_male_d, pcb_w/2 - conn_male_w/2, pcb_t + z_extra]) {
            difference() {
                union() {
                    // Outer Housing Block
                    cube([conn_male_d, conn_male_w, conn_male_h]);

                    // --- THE LATCH PROTRUSION (WEDGE PROFILE) ---
                    translate([conn_male_d, (conn_male_w/2) - 1.8, conn_male_h - 2.0]) {
                        hull() {
                            // Bottom flat face
                            cube([0.6, 3.6, 0.01]);
                            // Top edge tapering into the wall
                            translate([0, 0, 1.99]) {
                                cube([0.01, 3.6, 0.01]);
                            }
                        }
                    }
                }

                // Top-down hollow cavity to receive the female plug
                translate([0.8, 0.8, 1.5]) {
                    cube([conn_male_d - 1.6, conn_male_w - 1.6, conn_male_h]);
                }

                // --- TWO SPACINGS (SLOTS) ---
                for(i = [-1, 1]) {
                    translate([conn_male_d - 0.8, (conn_male_w/2) + (i * conn_pitch) - 0.7, 1.5]) {
                        cube([1.0, 1.4, conn_male_h]);
                    }
                }
            }
        }
    }

    // 2. Metallic Pins
    color("Silver") {
        pin_tip_h = 0.4; // cosmetic taper length of the solder point (unchanged)
        for(i = [-1, 0, 1]) {
            // Pointed solder tip protruding through the bottom of the PCB
            // (tip's sharp point lands at -pin_protrusion, derived from the
            // measured 12.115mm total stack: connector top -> pin tip end)
            translate([pcb_l - pin_offset_x, pcb_w/2 + (i * conn_pitch), -pin_protrusion + z_extra]) {
                cylinder(d1=0, d2=0.64, h=pin_tip_h, $fn=16);
            }
            // Main pin bodies extending up inside the shroud cavity
            translate([pcb_l - pin_offset_x, pcb_w/2 + (i * conn_pitch), -pin_protrusion + pin_tip_h + z_extra]) {
                cylinder(d=0.64, h=(conn_male_h - 1.5) + pcb_t + pin_protrusion - pin_tip_h, $fn=16);
            }
        }
    }
}

module connector_female(z_extra = 0) {
    translate([0, 0, z_extra]) {
        // 1. Female Plug Housing (with 3 hollow terminal cavities)
        color("Gainsboro", 0.98) {
            difference() {
                union() {
                    translate([pcb_l - conn_male_d + 0.85, pcb_w/2 - (conn_male_w - 1.7)/2, pcb_t + 1.5]) {

                        // Main plug block nested inside the male cavity
                        cube([conn_male_d - 1.7, conn_male_w - 1.7, 6.2]);

                        // Top structural flange/lip (The "blade" closest to the wires)
                        translate([-0.2, -0.4, 5.2]) {
                            cube([conn_male_d - 1.3, conn_male_w - 0.9, 1.2]);
                        }

                        // Two vertical guide ridges dropping down from the horizontal bar
                        for(i = [-1, 1]) {
                            translate([conn_male_d - 1.7, (conn_male_w - 1.7)/2 + (i * conn_pitch) - 0.6, 0]) {
                                cube([0.85, 1.2, 5.2]);
                            }
                        }

                        // Horizontal line connecting the 2 ridges
                        translate([conn_male_d - 1.7, (conn_male_w - 1.7)/2 - conn_pitch - 0.6, 5.2]) {
                            cube([0.85, (2 * conn_pitch) + 1.2, 1.2]);
                        }
                    }
                }

                // --- 3 RECTANGULAR TERMINAL HOLES (3.0mm x 2.0mm) ---
                // Subtracted down into the top face of the plug
                for(i = [-1, 0, 1]) {
                    translate([
                        (pcb_l - pin_offset_x) - 1.5,
                        (pcb_w / 2) + (i * conn_pitch) - 1.0,
                        pcb_t + 1.5 + 3.0 // Starts 3mm up inside the plug, punches through the top
                    ]) {
                        cube([3.0, 2.0, 4.0]);
                    }
                }
            }
        }

        // 2. Insulated Sensor Wires Emerging from Inside the Terminal Holes
        wire_colors = ["Red", "White", "Black"];
        for(i = [-1, 0, 1]) {
            color(wire_colors[i+1]) {
                translate([
                    pcb_l - pin_offset_x,
                    pcb_w/2 + (i * conn_pitch),
                    pcb_t + 1.5 + 3.0 // Anchored deep inside the hollow pocket
                ]) {
                    cylinder(d=1.1, h=20.0, $fn=16);
                }
            }
        }
    }
}

module safe_line_marker() {
    color("Black") {
        // Centered at safe_line_x, sitting on top of the PCB surface
        translate([safe_line_x, 0, pcb_t]) {
            hull() {
                // Bottom point: 2.5mm gap from edge + 0.75mm radius = Y at 3.25
                translate([0, 3.25, 0])
                    cylinder(d=1.5, h=0.05, $fn=32);

                // Top point: pcb_w - (2.5mm gap + 0.75mm radius) = Y at pcb_w - 3.25
                translate([0, pcb_w - 3.25, 0])
                    cylinder(d=1.5, h=0.05, $fn=32);
            }
        }
    }
}

module red_zone_markers() {
    color("Red") {
        // First red line (0.9cm / 9.0mm from right edge)
        translate([pcb_l - red_line_near_edge - 0.75, 1, pcb_t]) {
            cube([1.5, pcb_w - 2, 0.05]);
        }

        // Second red line (2.6cm / 26.0mm from right edge)
        translate([pcb_l - red_line_far_edge - 0.75, 1, pcb_t]) {
            cube([1.5, pcb_w - 2, 0.05]);
        }
    }
}

// =====================================================
// CONNECTOR FOOTPRINT OUTLINE (silkscreen / fit-check)
// =====================================================
// Draws a white rectangle on the top face of the PCB at
// exactly the position and size of the male shroud body
// (conn_male_d × conn_male_w), matching the translate()
// origin used in connector_male().
// Line thickness: nozzle_d (0.4 mm) — visually clear at
// 1:1 scale without overlapping the shroud body.
silk_t = nozzle_d;  // silkscreen line thickness (one nozzle width)
silk_z = 0.05;      // silk layer height above PCB surface

module connector_footprint_outline() {
    // Shroud origin on the PCB — mirrors connector_male() placement exactly
    ox = pcb_l - conn_male_d;
    oy = pcb_w / 2 - conn_male_w / 2;

    color("White") translate([ox, oy, pcb_t]) {
        // Bottom edge
        translate([0, 0, 0])
            cube([conn_male_d, silk_t, silk_z]);
        // Top edge
        translate([0, conn_male_w - silk_t, 0])
            cube([conn_male_d, silk_t, silk_z]);
        // Left edge
        translate([0, silk_t, 0])
            cube([silk_t, conn_male_w - 2 * silk_t, silk_z]);
        // Right edge
        translate([conn_male_d - silk_t, silk_t, 0])
            cube([silk_t, conn_male_w - 2 * silk_t, silk_z]);
    }
}

module pcb_assembly(male_z_extra = 0, female_z_extra = 0) {
    pcb_board();
    connector_male(z_extra = male_z_extra);
    connector_female(z_extra = female_z_extra);
    safe_line_marker();
    red_zone_markers();
    connector_footprint_outline();
}

// =====================================================
// RENDERING MODULES
// =====================================================
module full_system() {
    color("SteelBlue", 0.65) bottom_shell();

    translate([0, 0, z_pcb_seat]) pcb_assembly();
    translate([0, 0, outer_h]) color("LightBlue", 0.50) lid();
}

module exploded_system() {
    gap = 25;

    // Layer 0 — bottom shell (sits on the bed)
    color("SteelBlue", 0.85) bottom_shell();

    // Layer 1 — PCB board only (no connectors)
    translate([0, 0, z_pcb_seat + gap]) {
        pcb_board();
        safe_line_marker();
        red_zone_markers();
        connector_footprint_outline();
    }

    // Layer 2 — male connector + pins
    // Placed at the PCB seat height so its pcb_t-anchored geometry
    // is correct; z_extra lifts it an additional gap above the PCB.
    translate([0, 0, z_pcb_seat + gap]) {
        connector_male(z_extra = gap);
    }

    // Layer 3 — female plug + wires
    // Sits one further gap above the male shroud top (conn_male_h).
    translate([0, 0, z_pcb_seat + gap]) {
        connector_female(z_extra = gap + conn_male_h + gap);
    }

    // Layer 4 — lid
    translate([0, 0, outer_h + gap * 4]) color("LightBlue", 0.75) lid();
}

// =====================================================
// RENDER EXECUTION
// =====================================================
if (view_mode == 1) {
    full_system();
} else if (view_mode == 2) {
    exploded_system();
} else if (view_mode == 3) {
    difference() {
        full_system();
        translate([-10, pcb_w/2, -5]) {
            cube([pcb_l + 30, pcb_w + 20, outer_h + 30]);
        }
    }
}